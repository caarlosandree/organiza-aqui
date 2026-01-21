package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type RecurrenceHandler struct {
	recurrenceService service.RecurrenceService
	validator         *validator.Validate
}

func NewRecurrenceHandler(recurrenceService service.RecurrenceService) *RecurrenceHandler {
	return &RecurrenceHandler{
		recurrenceService: recurrenceService,
		validator:         validator.New(),
	}
}

// CreatePattern cria um novo padrão de recorrência
// @Summary Criar padrão de recorrência
// @Description Cria um padrão de recorrência para uma transação
// @Tags recurrence
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRecurrencePatternRequest true "Dados do padrão"
// @Success 201 {object} dto.RecurrencePatternDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/recurrence [post]
func (h *RecurrenceHandler) CreatePattern(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateRecurrencePatternRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	pattern, err := h.recurrenceService.CreatePattern(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar padrão de recorrência", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, pattern)
}

// GetPattern obtém um padrão por ID
// @Summary Obter padrão de recorrência
// @Description Retorna um padrão de recorrência específico
// @Tags recurrence
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do padrão"
// @Success 200 {object} dto.RecurrencePatternDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/recurrence/{id} [get]
func (h *RecurrenceHandler) GetPattern(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	patternID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	pattern, err := h.recurrenceService.GetPatternByID(c.Request().Context(), userID, patternID)
	if err != nil {
		if errors.Is(err, appError.ErrRecurrencePatternNotFound) {
			return response.NotFound(c, appError.ErrRecurrencePatternNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, pattern)
}

// GetPatterns lista todos os padrões do usuário
// @Summary Listar padrões de recorrência
// @Description Retorna todos os padrões de recorrência do usuário
// @Tags recurrence
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.RecurrencePatternDTO
// @Router /api/v1/recurrence [get]
func (h *RecurrenceHandler) GetPatterns(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	patterns, err := h.recurrenceService.GetPatternsByUserID(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar padrões", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, patterns)
}

// UpdatePattern atualiza um padrão
// @Summary Atualizar padrão de recorrência
// @Description Atualiza os dados de um padrão de recorrência
// @Tags recurrence
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do padrão"
// @Param request body dto.UpdateRecurrencePatternRequest true "Dados atualizados"
// @Success 200 {object} dto.RecurrencePatternDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/recurrence/{id} [put]
func (h *RecurrenceHandler) UpdatePattern(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	patternID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateRecurrencePatternRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	pattern, err := h.recurrenceService.UpdatePattern(c.Request().Context(), userID, patternID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrRecurrencePatternNotFound) {
			return response.NotFound(c, appError.ErrRecurrencePatternNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, pattern)
}

// DeletePattern deleta um padrão
// @Summary Deletar padrão de recorrência
// @Description Remove um padrão de recorrência
// @Tags recurrence
// @Security BearerAuth
// @Param id path string true "ID do padrão"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/recurrence/{id} [delete]
func (h *RecurrenceHandler) DeletePattern(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	patternID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.recurrenceService.DeletePattern(c.Request().Context(), userID, patternID); err != nil {
		if errors.Is(err, appError.ErrRecurrencePatternNotFound) {
			return response.NotFound(c, appError.ErrRecurrencePatternNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// GenerateTransactions gera transações recorrentes
// @Summary Gerar transações recorrentes
// @Description Gera transações baseadas nos padrões de recorrência até uma data
// @Tags recurrence
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GenerateRecurrenceRequest true "Data limite"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/recurrence/generate [post]
func (h *RecurrenceHandler) GenerateTransactions(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.GenerateRecurrenceRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	untilDate, err := time.Parse("2006-01-02", req.UntilDate)
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	count, err := h.recurrenceService.GenerateTransactions(c.Request().Context(), userID, untilDate)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao gerar transações", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, map[string]interface{}{
		"generated_count": count,
		"message":         fmt.Sprintf("%d transação(ões) gerada(s)", count),
	})
}

