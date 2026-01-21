package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/redis/go-redis/v9"
)

// UserRepository define a interface para operações de usuário
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error)
	GenerateUsernameSuggestions(ctx context.Context, baseUsername string, count int) ([]string, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

// NewUserRepository cria uma nova instância de UserRepository
func NewUserRepository(db *sqlx.DB, redisClient *redis.Client) UserRepository {
	return &userRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, name, created_at, updated_at)
		VALUES (:id, :email, :username, :password_hash, :name, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	// Tentar buscar do cache primeiro
	cacheKey := fmt.Sprintf("user:id:%s", id.String())
	cachedUserJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - deserializar e retornar
		var cachedUser model.User
		if err := json.Unmarshal([]byte(cachedUserJSON), &cachedUser); err == nil {
			return &cachedUser, nil
		}
	}

	// Cache miss - buscar no banco
	var user model.User
	query := `SELECT id, email, username, password_hash, name, created_at, updated_at FROM users WHERE id = $1`
	err = r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Cachear resultado (TTL: 30 minutos)
	userJSON, err := json.Marshal(user)
	if err == nil {
		// Se falhar ao cachear, não retornar erro (degradação graciosa)
		r.redis.Set(ctx, cacheKey, userJSON, 30*time.Minute)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	// Tentar buscar do cache primeiro
	cacheKey := fmt.Sprintf("user:email:%s", email)
	cachedUserJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - deserializar e retornar
		var cachedUser model.User
		if err := json.Unmarshal([]byte(cachedUserJSON), &cachedUser); err == nil {
			return &cachedUser, nil
		}
	}

	// Cache miss - buscar no banco
	var user model.User
	query := `SELECT id, email, username, password_hash, name, created_at, updated_at FROM users WHERE email = $1`
	err = r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Cachear resultado (TTL: 30 minutos)
	userJSON, err := json.Marshal(user)
	if err == nil {
		// Cachear por ID também
		idCacheKey := fmt.Sprintf("user:id:%s", user.ID.String())
		r.redis.Set(ctx, idCacheKey, userJSON, 30*time.Minute)
		// Cachear por email
		r.redis.Set(ctx, cacheKey, userJSON, 30*time.Minute)
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	// Tentar buscar do cache primeiro
	cacheKey := fmt.Sprintf("user:username:%s", username)
	cachedUserJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - deserializar e retornar
		var cachedUser model.User
		if err := json.Unmarshal([]byte(cachedUserJSON), &cachedUser); err == nil {
			return &cachedUser, nil
		}
	}

	// Cache miss - buscar no banco
	var user model.User
	query := `SELECT id, email, username, password_hash, name, created_at, updated_at FROM users WHERE username = $1`
	err = r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Cachear resultado (TTL: 30 minutos)
	userJSON, err := json.Marshal(user)
	if err == nil {
		// Cachear por ID, email e username
		idCacheKey := fmt.Sprintf("user:id:%s", user.ID.String())
		emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
		r.redis.Set(ctx, idCacheKey, userJSON, 30*time.Minute)
		r.redis.Set(ctx, emailCacheKey, userJSON, 30*time.Minute)
		r.redis.Set(ctx, cacheKey, userJSON, 30*time.Minute)
	}

	return &user, nil
}

func (r *userRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	// Primeiro tenta como email, depois como username
	user, err := r.FindByEmail(ctx, identifier)
	if err != nil || user != nil {
		return user, err
	}
	return r.FindByUsername(ctx, identifier)
}

func (r *userRepository) GenerateUsernameSuggestions(ctx context.Context, baseUsername string, count int) ([]string, error) {
	suggestions := make([]string, 0, count)
	
	// Gerar sugestões adicionando números sequenciais
	for i := 1; len(suggestions) < count && i <= count*2; i++ {
		suggestion := fmt.Sprintf("%s%d", baseUsername, i)
		
		// Verificar se a sugestão já existe
		existing, err := r.FindByUsername(ctx, suggestion)
		if err != nil {
			return nil, err
		}
		
		// Se não existe, adicionar à lista
		if existing == nil {
			suggestions = append(suggestions, suggestion)
		}
	}
	
	// Se ainda não temos sugestões suficientes, tentar com sufixos aleatórios
	if len(suggestions) < count {
		attempts := 0
		maxAttempts := (count - len(suggestions)) * 10
		
		for len(suggestions) < count && attempts < maxAttempts {
			attempts++
			// Gerar sufixo aleatório de 3-4 dígitos
			suffix := fmt.Sprintf("%d", time.Now().UnixNano()%10000)
			suggestion := fmt.Sprintf("%s%s", baseUsername, suffix)
			
			existing, err := r.FindByUsername(ctx, suggestion)
			if err != nil {
				continue
			}
			
			if existing == nil {
				// Verificar se já não está na lista
				exists := false
				for _, s := range suggestions {
					if s == suggestion {
						exists = true
						break
					}
				}
				if !exists {
					suggestions = append(suggestions, suggestion)
				}
			}
		}
	}
	
	return suggestions, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Buscar usuário antigo para invalidar cache por email
	var oldEmail string
	err := r.db.GetContext(ctx, &oldEmail, `SELECT email FROM users WHERE id = $1`, user.ID)
	if err != nil && err != sql.ErrNoRows {
		// Se houver erro (exceto não encontrado), continuar mesmo assim
	}

	query := `
		UPDATE users 
		SET email = :email, username = :username, password_hash = :password_hash, name = :name, updated_at = :updated_at
		WHERE id = :id
	`
	_, err = r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	// Buscar username antigo para invalidar cache
	var oldUsername string
	r.db.GetContext(ctx, &oldUsername, `SELECT username FROM users WHERE id = $1`, user.ID)
	
	// Invalidar cache
	r.invalidateUserCache(ctx, user.ID, user.Email, user.Username)
	if oldEmail != "" && oldEmail != user.Email {
		// Email mudou, invalidar cache antigo também
		oldCacheKey := fmt.Sprintf("user:email:%s", oldEmail)
		r.redis.Del(ctx, oldCacheKey)
	}
	if oldUsername != "" && oldUsername != user.Username {
		// Username mudou, invalidar cache antigo também
		oldCacheKey := fmt.Sprintf("user:username:%s", oldUsername)
		r.redis.Del(ctx, oldCacheKey)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Buscar usuário antes de deletar para invalidar cache
	user, _ := r.FindByID(ctx, id)

	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Invalidar cache
	if user != nil {
		r.invalidateUserCache(ctx, user.ID, user.Email, user.Username)
	}

	return nil
}

// invalidateUserCache remove o cache do usuário (por ID, email e username)
func (r *userRepository) invalidateUserCache(ctx context.Context, id uuid.UUID, email string, username string) {
	idCacheKey := fmt.Sprintf("user:id:%s", id.String())
	emailCacheKey := fmt.Sprintf("user:email:%s", email)
	usernameCacheKey := fmt.Sprintf("user:username:%s", username)
	r.redis.Del(ctx, idCacheKey, emailCacheKey, usernameCacheKey)
}
