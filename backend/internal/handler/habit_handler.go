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

type HabitHandler struct {
	habitService service.HabitService
	validator    *validator.Validate
}

func NewHabitHandler(habitService service.HabitService) *HabitHandler {
	return &HabitHandler{
		habitService: habitService,
		validator:    validator.New(),
	}
}

// CreateHabit cria um novo hábito
// @Summary Criar hábito
// @Description Cria um novo hábito
// @Tags habits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateHabitRequest true "Dados do hábito"
// @Success 201 {object} dto.HabitDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/habits [post]
func (h *HabitHandler) CreateHabit(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateHabitRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	habit, err := h.habitService.CreateHabit(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar hábito", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusCreated, habit)
}

// GetHabit obtém um hábito por ID
// @Summary Obter hábito
// @Description Retorna um hábito específico
// @Tags habits
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do hábito"
// @Success 200 {object} dto.HabitDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/habits/{id} [get]
func (h *HabitHandler) GetHabit(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	habitID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	habit, err := h.habitService.GetHabitByID(c.Request().Context(), userID, habitID)
	if err != nil {
		if errors.Is(err, appError.ErrHabitNotFound) {
			return response.NotFound(c, appError.ErrHabitNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, habit)
}

// GetHabits lista todos os hábitos do usuário
// @Summary Listar hábitos
// @Description Retorna todos os hábitos do usuário
// @Tags habits
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.HabitDTO
// @Router /api/v1/habits [get]
func (h *HabitHandler) GetHabits(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	habits, err := h.habitService.GetHabitsByUserID(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar hábitos", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, habits)
}

// UpdateHabit atualiza um hábito
// @Summary Atualizar hábito
// @Description Atualiza os dados de um hábito
// @Tags habits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do hábito"
// @Param request body dto.UpdateHabitRequest true "Dados atualizados"
// @Success 200 {object} dto.HabitDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/habits/{id} [put]
func (h *HabitHandler) UpdateHabit(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	habitID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateHabitRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	habit, err := h.habitService.UpdateHabit(c.Request().Context(), userID, habitID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrHabitNotFound) {
			return response.NotFound(c, appError.ErrHabitNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, habit)
}

// DeleteHabit deleta um hábito
// @Summary Deletar hábito
// @Description Remove um hábito
// @Tags habits
// @Security BearerAuth
// @Param id path string true "ID do hábito"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/habits/{id} [delete]
func (h *HabitHandler) DeleteHabit(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	habitID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.habitService.DeleteHabit(c.Request().Context(), userID, habitID); err != nil {
		if errors.Is(err, appError.ErrHabitNotFound) {
			return response.NotFound(c, appError.ErrHabitNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetHabitStats obtém estatísticas de um hábito
// @Summary Obter estatísticas do hábito
// @Description Retorna estatísticas de um hábito (completion rate, streaks, etc.)
// @Tags habits
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do hábito"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Success 200 {object} dto.HabitStatsDTO
// @Router /api/v1/habits/{id}/stats [get]
func (h *HabitHandler) GetHabitStats(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	habitID, err := parseUUIDParam(c, "id")
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

	stats, err := h.habitService.GetHabitStats(c.Request().Context(), userID, habitID, startDate, endDate)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar estatísticas", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, stats)
}
