package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type CalendarEventService interface {
	CreateEvent(ctx context.Context, userID uuid.UUID, req *dto.CreateCalendarEventRequest) (*dto.CalendarEventDTO, error)
	GetEventByID(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.CalendarEventDTO, error)
	GetEventsByUserID(ctx context.Context, userID uuid.UUID, filters *dto.CalendarEventFilters) ([]*dto.CalendarEventDTO, error)
	UpdateEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req *dto.UpdateCalendarEventRequest) (*dto.CalendarEventDTO, error)
	DeleteEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error
}

type calendarEventService struct {
	eventRepo repository.CalendarEventRepository
}

func NewCalendarEventService(eventRepo repository.CalendarEventRepository) CalendarEventService {
	return &calendarEventService{
		eventRepo: eventRepo,
	}
}

func (s *calendarEventService) CreateEvent(ctx context.Context, userID uuid.UUID, req *dto.CreateCalendarEventRequest) (*dto.CalendarEventDTO, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("start_date inválida: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date inválida: %w", err)
		}
		endDate = &parsed
	}

	now := time.Now()
	event := &model.CalendarEvent{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		AllDay:      req.AllDay,
		Location:    req.Location,
		Color:       req.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("erro ao criar evento: %w", err)
	}

	return s.modelToDTO(event), nil
}

func (s *calendarEventService) GetEventByID(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) (*dto.CalendarEventDTO, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar evento: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("evento %s: %w", eventID, appError.ErrCalendarEventNotFound)
	}
	if event.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(event), nil
}

func (s *calendarEventService) GetEventsByUserID(ctx context.Context, userID uuid.UUID, filters *dto.CalendarEventFilters) ([]*dto.CalendarEventDTO, error) {
	events, err := s.eventRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos: %w", err)
	}

	dtos := make([]*dto.CalendarEventDTO, len(events))
	for i, event := range events {
		dtos[i] = s.modelToDTO(event)
	}

	return dtos, nil
}

func (s *calendarEventService) UpdateEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, req *dto.UpdateCalendarEventRequest) (*dto.CalendarEventDTO, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar evento: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("evento %s: %w", eventID, appError.ErrCalendarEventNotFound)
	}
	if event.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("start_date inválida: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date inválida: %w", err)
		}
		endDate = &parsed
	}

	event.Title = req.Title
	event.Description = req.Description
	event.StartDate = startDate
	event.EndDate = endDate
	event.AllDay = req.AllDay
	event.Location = req.Location
	event.Color = req.Color
	event.UpdatedAt = time.Now()

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("erro ao atualizar evento: %w", err)
	}

	return s.modelToDTO(event), nil
}

func (s *calendarEventService) DeleteEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("erro ao buscar evento: %w", err)
	}
	if event == nil {
		return fmt.Errorf("evento %s: %w", eventID, appError.ErrCalendarEventNotFound)
	}
	if event.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.eventRepo.Delete(ctx, eventID); err != nil {
		return fmt.Errorf("erro ao deletar evento: %w", err)
	}

	return nil
}

func (s *calendarEventService) modelToDTO(event *model.CalendarEvent) *dto.CalendarEventDTO {
	dto := &dto.CalendarEventDTO{
		ID:          event.ID.String(),
		UserID:      event.UserID.String(),
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate.Format(time.RFC3339),
		AllDay:      event.AllDay,
		Location:    event.Location,
		Color:       event.Color,
		CreatedAt:   event.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   event.UpdatedAt.Format(time.RFC3339),
	}

	if event.EndDate != nil {
		endDateStr := event.EndDate.Format(time.RFC3339)
		dto.EndDate = &endDateStr
	}

	return dto
}
