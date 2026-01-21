package response

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	appError "github.com/organiza-aqui/backend/internal/error"
)

// ErrorResponse representa uma resposta de erro padrão
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse representa uma resposta de sucesso padrão
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON retorna uma resposta JSON padrão
func JSON(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, data)
}

// Success retorna uma resposta de sucesso
func Success(c echo.Context, message string, data interface{}) error {
	return JSON(c, http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// Error retorna uma resposta de erro
func Error(c echo.Context, statusCode int, err error) error {
	if err == nil {
		err = errors.New(http.StatusText(statusCode))
	}
	return JSON(c, statusCode, ErrorResponse{
		Error:   err.Error(),
		Message: http.StatusText(statusCode),
	})
}

// ErrorWithMessage retorna uma resposta de erro com mensagem customizada
func ErrorWithMessage(c echo.Context, statusCode int, err error, message string) error {
	if err == nil {
		err = errors.New(http.StatusText(statusCode))
	}
	return JSON(c, statusCode, ErrorResponse{
		Error:   err.Error(),
		Message: message,
	})
}

// BadRequest retorna um erro 400
func BadRequest(c echo.Context, err error) error {
	return Error(c, http.StatusBadRequest, err)
}

// Unauthorized retorna um erro 401
func Unauthorized(c echo.Context, err error) error {
	return Error(c, http.StatusUnauthorized, err)
}

// NotFound retorna um erro 404
func NotFound(c echo.Context, err error) error {
	return Error(c, http.StatusNotFound, err)
}

// Conflict retorna um erro 409
func Conflict(c echo.Context, err error) error {
	return Error(c, http.StatusConflict, err)
}

// InternalServerError retorna um erro 500
// Sempre usa ErrInternalServer para evitar expor detalhes internos
func InternalServerError(c echo.Context, err error) error {
	// Sempre usar o erro sanitizado para não expor detalhes internos
	return Error(c, http.StatusInternalServerError, appError.ErrInternalServer)
}