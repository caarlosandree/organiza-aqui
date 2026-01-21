package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/response"
)

// AuthMiddleware valida o token JWT e adiciona o usuário ao contexto
func AuthMiddleware(authService service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Obter token do header Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, echo.NewHTTPError(http.StatusUnauthorized, "token de autenticação não fornecido"))
			}

			// Extrair token (formato: "Bearer <token>")
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Unauthorized(c, echo.NewHTTPError(http.StatusUnauthorized, "formato de token inválido"))
			}

			token := parts[1]

			// Validar token
			user, err := authService.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return response.Unauthorized(c, err)
			}

			// Adicionar usuário ao contexto
			c.Set("user", user)
			c.Set("user_id", user.ID)

			return next(c)
		}
	}
}
