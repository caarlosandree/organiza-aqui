package handler

import (
	"errors"
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

type HabitTrackingHandler struct {
	trackingService service.HabitTrackingService
	validator       *validator.Validate
}

func NewHabitTrackingHandler(trackingService service.HabitTrackingService) *HabitTrackingHandler {
	return &HabitTrackingHandler{
		trackingService: trackingService,
		validator:       validator.New(),
	}
}

// CreateTracking cria um novo registro de tracking
// @Summary Criar registro de tracking
// @Description Cria um novo registro de execução de hábito
// @Tags habit-tracking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateHabitTrackingRequest true "Dados do tracking"
// @Success 201 {object} dto.HabitTrackingDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/habit-tracking [post]
func (h *HabitTrackingHandler) CreateTracking(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateHabitTrackingRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	tracking, err := h.trackingService.CreateTracking(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar tracking", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, tracking)
}

// GetTracking obtém um tracking por ID
// @Summary Obter tracking
// @Description Retorna um tracking específico
// @Tags habit-tracking
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do tracking"
// @Success 200 {object} dto.HabitTrackingDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/habit-tracking/{id} [get]
func (h *HabitTrackingHandler) GetTracking(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	trackingID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	tracking, err := h.trackingService.GetTrackingByID(c.Request().Context(), userID, trackingID)
	if err != nil {
		if errors.Is(err, appError.ErrTrackingNotFound) {
			return response.NotFound(c, appError.ErrTrackingNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, tracking)
}

// GetTrackingByHabit lista todos os trackings de um hábito
// @Summary Listar trackings de um hábito
// @Description Retorna todos os trackings de um hábito em um período
// @Tags habit-tracking
// @Produce json
// @Security BearerAuth
// @Param habit_id path string true "ID do hábito"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Success 200 {array} dto.HabitTrackingDTO
// @Router /api/v1/habits/{habit_id}/tracking [get]
func (h *HabitTrackingHandler) GetTrackingByHabit(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	habitID, err := parseUUIDParam(c, "habit_id")
	if err != nil {
		return err
	}

	// Default: últimos 30 dias
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	trackings, err := h.trackingService.GetTrackingByHabitID(c.Request().Context(), userID, habitID, startDate, endDate)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar trackings", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, trackings)
}

// UpdateTracking atualiza um tracking
// @Summary Atualizar tracking
// @Description Atualiza os dados de um tracking
// @Tags habit-tracking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do tracking"
// @Param request body dto.UpdateHabitTrackingRequest true "Dados atualizados"
// @Success 200 {object} dto.HabitTrackingDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/habit-tracking/{id} [put]
func (h *HabitTrackingHandler) UpdateTracking(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	trackingID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateHabitTrackingRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	tracking, err := h.trackingService.UpdateTracking(c.Request().Context(), userID, trackingID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrTrackingNotFound) {
			return response.NotFound(c, appError.ErrTrackingNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, tracking)
}

// DeleteTracking deleta um tracking
// @Summary Deletar tracking
// @Description Remove um registro de tracking
// @Tags habit-tracking
// @Security BearerAuth
// @Param id path string true "ID do tracking"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/habit-tracking/{id} [delete]
func (h *HabitTrackingHandler) DeleteTracking(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	trackingID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.trackingService.DeleteTracking(c.Request().Context(), userID, trackingID); err != nil {
		if errors.Is(err, appError.ErrTrackingNotFound) {
			return response.NotFound(c, appError.ErrTrackingNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

