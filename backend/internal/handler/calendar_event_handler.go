package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type CalendarEventHandler struct {
	eventService service.CalendarEventService
	validator    *validator.Validate
}

func NewCalendarEventHandler(eventService service.CalendarEventService) *CalendarEventHandler {
	return &CalendarEventHandler{
		eventService: eventService,
		validator:    validator.New(),
	}
}

// CreateEvent cria um novo evento
// @Summary Criar evento do calendário
// @Description Cria um novo evento no calendário
// @Tags calendar-events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCalendarEventRequest true "Dados do evento"
// @Success 201 {object} dto.CalendarEventDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/calendar-events [post]
func (h *CalendarEventHandler) CreateEvent(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateCalendarEventRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	event, err := h.eventService.CreateEvent(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar evento", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, event)
}

// GetEvent obtém um evento por ID
// @Summary Obter evento do calendário
// @Description Retorna um evento específico
// @Tags calendar-events
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do evento"
// @Success 200 {object} dto.CalendarEventDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/calendar-events/{id} [get]
func (h *CalendarEventHandler) GetEvent(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	event, err := h.eventService.GetEventByID(c.Request().Context(), userID, eventID)
	if err != nil {
		if errors.Is(err, appError.ErrCalendarEventNotFound) {
			return response.NotFound(c, appError.ErrCalendarEventNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, event)
}

// GetEvents lista todos os eventos do usuário
// @Summary Listar eventos do calendário
// @Description Retorna todos os eventos do usuário com filtros opcionais
// @Tags calendar-events
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Data inicial (RFC3339)"
// @Param end_date query string false "Data final (RFC3339)"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {array} dto.CalendarEventDTO
// @Router /api/v1/calendar-events [get]
func (h *CalendarEventHandler) GetEvents(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	filters := &dto.CalendarEventFilters{}

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

	events, err := h.eventService.GetEventsByUserID(c.Request().Context(), userID, filters)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar eventos", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, events)
}

// UpdateEvent atualiza um evento
// @Summary Atualizar evento do calendário
// @Description Atualiza os dados de um evento
// @Tags calendar-events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do evento"
// @Param request body dto.UpdateCalendarEventRequest true "Dados atualizados"
// @Success 200 {object} dto.CalendarEventDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/calendar-events/{id} [put]
func (h *CalendarEventHandler) UpdateEvent(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateCalendarEventRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	event, err := h.eventService.UpdateEvent(c.Request().Context(), userID, eventID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrCalendarEventNotFound) {
			return response.NotFound(c, appError.ErrCalendarEventNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, event)
}

// DeleteEvent deleta um evento
// @Summary Deletar evento do calendário
// @Description Remove um evento do calendário
// @Tags calendar-events
// @Security BearerAuth
// @Param id path string true "ID do evento"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/calendar-events/{id} [delete]
func (h *CalendarEventHandler) DeleteEvent(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	eventID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.eventService.DeleteEvent(c.Request().Context(), userID, eventID); err != nil {
		if errors.Is(err, appError.ErrCalendarEventNotFound) {
			return response.NotFound(c, appError.ErrCalendarEventNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

