package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/redis/go-redis/v9"

	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

// AuthService define a interface para operações de autenticação
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*dto.UserDTO, error)
	Logout(ctx context.Context, tokenString string) error
	GenerateUsernameSuggestions(ctx context.Context, baseUsername string) ([]string, error)
}

type authService struct {
	userRepo           repository.UserRepository
	authRepo           repository.AuthRepository
	redis              *redis.Client
	jwtSecret          string
	jwtExpirationHours  int
	cacheTokenTTL      time.Duration
}

// NewAuthService cria uma nova instância de AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	redisClient *redis.Client,
	jwtSecret string,
	jwtExpirationHours int,
) AuthService {
	cacheTokenTTL := 5 * time.Minute // TTL padrão de 5 minutos para cache de token
	
	return &authService{
		userRepo:           userRepo,
		authRepo:           authRepo,
		redis:              redisClient,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
		cacheTokenTTL:      cacheTokenTTL,
	}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Validar força da senha (validação adicional de segurança)
	if !validatePasswordStrength(req.Password) {
		return nil, errors.New("senha não atende aos critérios de força mínima (use letras maiúsculas, minúsculas, números e caracteres especiais)")
	}

	// Verificar se o email já existe
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar email: %w", err)
	}
	if existingUser != nil {
		return nil, appError.ErrEmailAlreadyExists
	}

	// Verificar se o username já existe
	existingUserByUsername, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar username: %w", err)
	}
	if existingUserByUsername != nil {
		return nil, appError.ErrUsernameAlreadyExists
	}

	// Hash da senha
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Criar usuário
	now := time.Now()
	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Name:         req.Name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// Gerar token JWT
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	// Criar sessão
	if err := s.createSession(ctx, user.ID, token); err != nil {
		return nil, fmt.Errorf("erro ao criar sessão: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
			Name:     user.Name,
		},
	}, nil
}

// GenerateUsernameSuggestions gera sugestões de username quando o username desejado já está em uso
func (s *authService) GenerateUsernameSuggestions(ctx context.Context, baseUsername string) ([]string, error) {
	return s.userRepo.GenerateUsernameSuggestions(ctx, baseUsername, 5)
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Buscar usuário por email ou username
	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.Identifier)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	if user == nil {
		return nil, appError.ErrInvalidCredentials
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, appError.ErrInvalidCredentials
	}

	// Gerar token JWT
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	// Criar sessão
	if err := s.createSession(ctx, user.ID, token); err != nil {
		return nil, fmt.Errorf("erro ao criar sessão: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
			Name:     user.Name,
		},
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*dto.UserDTO, error) {
	tokenHash := s.hashToken(tokenString)
	cacheKey := fmt.Sprintf("token:valid:%s", tokenHash)

	// Tentar buscar do cache primeiro
	cachedUserJSON, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - deserializar e retornar
		var cachedUser dto.UserDTO
		if err := json.Unmarshal([]byte(cachedUserJSON), &cachedUser); err == nil {
			return &cachedUser, nil
		}
	}

	// Cache miss ou erro - validar token normalmente
	// Parse do token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Extrair claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("erro ao extrair claims do token")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id não encontrado no token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("user_id inválido: %w", err)
	}

	// Verificar se a sessão existe e está válida
	session, err := s.authRepo.FindSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar sessão: %w", err)
	}
	if session == nil || session.UserID != userID {
		return nil, errors.New("sessão inválida ou expirada")
	}

	// Buscar usuário
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	userDTO := &dto.UserDTO{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		Name:     user.Name,
	}

	// Cachear resultado para próximas validações
	userJSON, err := json.Marshal(userDTO)
	if err == nil {
		// Se falhar ao cachear, não retornar erro (degradação graciosa)
		s.redis.Set(ctx, cacheKey, userJSON, s.cacheTokenTTL)
	}

	return userDTO, nil
}

func (s *authService) Logout(ctx context.Context, tokenString string) error {
	tokenHash := s.hashToken(tokenString)
	
	// Invalidar cache de token
	cacheKey := fmt.Sprintf("token:valid:%s", tokenHash)
	s.redis.Del(ctx, cacheKey)
	
	session, err := s.authRepo.FindSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("erro ao buscar sessão: %w", err)
	}
	if session != nil {
		if err := s.authRepo.DeleteSession(ctx, session.ID); err != nil {
			return fmt.Errorf("erro ao deletar sessão: %w", err)
		}
	}
	return nil
}

func (s *authService) generateJWT(userID uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) createSession(ctx context.Context, userID uuid.UUID, tokenString string) error {
	tokenHash := s.hashToken(tokenString)
	expiresAt := time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour)

	session := &model.AuthSession{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	return s.authRepo.CreateSession(ctx, session)
}

func (s *authService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// validatePasswordStrength valida a força da senha
// Retorna true se a senha atender aos critérios mínimos de força
func validatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Senha deve ter pelo menos 3 dos 4 critérios
	criteriaCount := 0
	if hasUpper {
		criteriaCount++
	}
	if hasLower {
		criteriaCount++
	}
	if hasNumber {
		criteriaCount++
	}
	if hasSpecial {
		criteriaCount++
	}

	return criteriaCount >= 3
}
