package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

// getUserID extrai o userID do contexto
func getUserID(c echo.Context) (uuid.UUID, error) {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return uuid.Nil, response.Unauthorized(c, appError.ErrUnauthorized)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, response.Unauthorized(c, appError.ErrUnauthorized)
	}

	return userID, nil
}

// parseUUIDParam extrai e valida um UUID de um parâmetro da rota
func parseUUIDParam(c echo.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, response.BadRequest(c, appError.ErrInvalidInput)
	}
	return id, nil
}

// bindAndValidate faz bind do request e valida os dados
// Retorna nil se houver erro (já enviou resposta HTTP)
func bindAndValidate(c echo.Context, req interface{}, validator *validator.Validate) error {
	if err := c.Bind(req); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao fazer bind do request", err)
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := validator.Struct(req); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro de validação", err)
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	return nil
}

// handleServiceError trata erros do service de forma padronizada
// Retorna nil se houver erro (já enviou resposta HTTP)
func handleServiceError(c echo.Context, err error, logMessage string) error {
	if err == nil {
		return nil
	}

	logger.ErrorWithContext(c.Request().Context(), logMessage, err)

	if errors.Is(err, appError.ErrNotFound) {
		return response.NotFound(c, appError.ErrNotFound)
	}

	if errors.Is(err, appError.ErrInvalidInput) {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if errors.Is(err, appError.ErrUnauthorizedAccess) {
		return response.Unauthorized(c, appError.ErrUnauthorized)
	}

	if errors.Is(err, appError.ErrConflict) {
		return response.Conflict(c, appError.ErrConflict)
	}

	// Erro genérico - retornar erro interno
	return response.InternalServerError(c, appError.ErrInternalServer)
}