package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type TimelineHandler struct {
	timelineService service.TimelineService
}

func NewTimelineHandler(timelineService service.TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
	}
}

// GetTimelineEvents obtém eventos da timeline
// @Summary Obter eventos da timeline
// @Description Retorna eventos da timeline unificada
// @Tags timeline
// @Produce json
// @Security BearerAuth
// @Param entity_type query string false "Filtrar por tipo de entidade"
// @Param start_date query string false "Data inicial"
// @Param end_date query string false "Data final"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {array} dto.TimelineEventDTO
// @Router /api/v1/timeline [get]
func (h *TimelineHandler) GetTimelineEvents(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	filters := &dto.TimelineFilters{}

	if entityType := c.QueryParam("entity_type"); entityType != "" {
		filters.EntityType = &entityType
	}

	if startDate := c.QueryParam("start_date"); startDate != "" {
		filters.StartDate = &startDate
	}

	if endDate := c.QueryParam("end_date"); endDate != "" {
		filters.EndDate = &endDate
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	events, err := h.timelineService.GetTimelineEvents(c.Request().Context(), userID, filters)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar eventos da timeline", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, events)
}

// GetTimelineSummary obtém resumo da timeline
// @Summary Obter resumo da timeline
// @Description Retorna um resumo da timeline (totais, eventos de hoje, futuros, por tipo)
// @Tags timeline
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.TimelineSummaryDTO
// @Router /api/v1/timeline/summary [get]
func (h *TimelineHandler) GetTimelineSummary(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	summary, err := h.timelineService.GetTimelineSummary(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar resumo da timeline", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, summary)
}

