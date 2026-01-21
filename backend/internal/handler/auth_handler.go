package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	passwordValidator "github.com/organiza-aqui/backend/internal/validator"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
	"go.uber.org/zap"
)

// AuthHandler gerencia handlers de autenticação
type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler cria uma nova instância de AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	v := validator.New()
	
	// Registrar validação customizada de força de senha
	v.RegisterValidation("password_strength", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return passwordValidator.ValidatePasswordStrength(password)
	})
	
	return &AuthHandler{
		authService: authService,
		validator:   v,
	}
}

// GetAuthService retorna o authService (para uso no middleware)
func (h *AuthHandler) GetAuthService() service.AuthService {
	return h.authService
}

// Register registra um novo usuário
// @Summary Registrar novo usuário
// @Description Cria uma nova conta de usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Dados de registro"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao fazer bind do request", err)
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro de validação", err)
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	resp, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao registrar usuário", err, zap.String("email", req.Email), zap.String("username", req.Username))
		
		// Verificar se é erro de email duplicado
		if errors.Is(err, appError.ErrEmailAlreadyExists) {
			return response.Conflict(c, appError.ErrEmailAlreadyExists)
		}

		// Verificar se é erro de username duplicado
		if errors.Is(err, appError.ErrUsernameAlreadyExists) {
			// Gerar sugestões de username
			suggestions, suggestErr := h.authService.GenerateUsernameSuggestions(c.Request().Context(), req.Username)
			if suggestErr != nil {
				logger.ErrorWithContext(c.Request().Context(), "Erro ao gerar sugestões de username", suggestErr)
			}

			// Criar resposta com sugestões
			errorData := map[string]interface{}{
				"error":       appError.ErrUsernameAlreadyExists.Error(),
				"message":     "este username já está em uso",
				"suggestions": suggestions,
			}

			return response.JSON(c, http.StatusConflict, errorData)
		}

		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusCreated, resp)
}

// Login autentica um usuário
// @Summary Login
// @Description Autentica um usuário e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciais de login"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao fazer bind do request", err)
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro de validação", err)
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	resp, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao fazer login", err, zap.String("identifier", req.Identifier))
		
		// Verificar se é erro de credenciais inválidas
		if errors.Is(err, appError.ErrInvalidCredentials) {
			return response.Unauthorized(c, appError.ErrInvalidCredentials)
		}

		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, resp)
}

// Logout encerra a sessão do usuário
// @Summary Logout
// @Description Encerra a sessão do usuário
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	// Obter token do header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return response.Unauthorized(c, appError.ErrUnauthorized)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return response.Unauthorized(c, appError.ErrUnauthorized)
	}

	token := parts[1]

	if err := h.authService.Logout(c.Request().Context(), token); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao fazer logout", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.Success(c, "logout realizado com sucesso", nil)
}

// Me retorna informações do usuário autenticado
// @Summary Obter usuário autenticado
// @Description Retorna informações do usuário autenticado
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserDTO
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	user, ok := c.Get("user").(*dto.UserDTO)
	if !ok || user == nil {
		return response.Unauthorized(c, appError.ErrUnauthorized)
	}

	return response.JSON(c, http.StatusOK, user)
}
