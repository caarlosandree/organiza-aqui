package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type TimelineService interface {
	GetTimelineEvents(ctx context.Context, userID uuid.UUID, filters *dto.TimelineFilters) ([]*dto.TimelineEventDTO, error)
	GetTimelineSummary(ctx context.Context, userID uuid.UUID) (*dto.TimelineSummaryDTO, error)
	SyncTimelineFromTransactions(ctx context.Context, userID uuid.UUID) error
	SyncTimelineFromTasks(ctx context.Context, userID uuid.UUID) error
}

type timelineService struct {
	timelineRepo     repository.TimelineRepository
	transactionRepo  repository.TransactionRepository
	taskRepo         repository.TaskRepository
}

func NewTimelineService(
	timelineRepo repository.TimelineRepository,
	transactionRepo repository.TransactionRepository,
	taskRepo repository.TaskRepository,
) TimelineService {
	return &timelineService{
		timelineRepo:    timelineRepo,
		transactionRepo: transactionRepo,
		taskRepo:        taskRepo,
	}
}

func (s *timelineService) GetTimelineEvents(ctx context.Context, userID uuid.UUID, filters *dto.TimelineFilters) ([]*dto.TimelineEventDTO, error) {
	events, err := s.timelineRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos: %w", err)
	}

	dtos := make([]*dto.TimelineEventDTO, len(events))
	for i, event := range events {
		dtos[i] = s.modelToDTO(event)
	}

	return dtos, nil
}

func (s *timelineService) GetTimelineSummary(ctx context.Context, userID uuid.UUID) (*dto.TimelineSummaryDTO, error) {
	summary, err := s.timelineRepo.GetSummary(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar resumo: %w", err)
	}

	return summary, nil
}

func (s *timelineService) SyncTimelineFromTransactions(ctx context.Context, userID uuid.UUID) error {
	// Buscar transações recentes (últimos 30 dias)
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	filters := &dto.TransactionFilters{
		StartDate: &startDate,
	}

	_, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return fmt.Errorf("erro ao buscar transações: %w", err)
	}

	// Criar eventos de timeline para cada transação
	// Nota: Em produção, isso seria feito via triggers ou jobs
	// Por enquanto, vamos apenas retornar sucesso
	// A implementação completa criaria registros na tabela timeline_events

	return nil
}

func (s *timelineService) SyncTimelineFromTasks(ctx context.Context, userID uuid.UUID) error {
	// Buscar tarefas recentes
	filters := &dto.TaskFilters{
		Limit: 100,
	}

	_, err := s.taskRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return fmt.Errorf("erro ao buscar tarefas: %w", err)
	}

	// Criar eventos de timeline para cada tarefa
	// Nota: Em produção, isso seria feito via triggers ou jobs
	// Por enquanto, vamos apenas retornar sucesso

	return nil
}

func (s *timelineService) modelToDTO(event *model.TimelineEvent) *dto.TimelineEventDTO {
	dto := &dto.TimelineEventDTO{
		ID:          event.ID.String(),
		UserID:      event.UserID.String(),
		EntityType:  event.EntityType,
		EntityID:    event.EntityID.String(),
		Title:       event.Title,
		Description: event.Description,
		EventDate:   event.EventDate.Format(time.RFC3339),
		CreatedAt:   event.CreatedAt.Format(time.RFC3339),
	}

	// Parse metadata JSON
	if len(event.Metadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(event.Metadata, &metadata); err == nil {
			dto.Metadata = metadata
		}
	}

	return dto
}
