package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type RecurrenceService interface {
	CreatePattern(ctx context.Context, userID uuid.UUID, req *dto.CreateRecurrencePatternRequest) (*dto.RecurrencePatternDTO, error)
	GetPatternByID(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) (*dto.RecurrencePatternDTO, error)
	GetPatternsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.RecurrencePatternDTO, error)
	UpdatePattern(ctx context.Context, userID uuid.UUID, patternID uuid.UUID, req *dto.UpdateRecurrencePatternRequest) (*dto.RecurrencePatternDTO, error)
	DeletePattern(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) error
	GenerateTransactions(ctx context.Context, userID uuid.UUID, untilDate time.Time) (int, error)
}

type recurrenceService struct {
	recurrenceRepo repository.RecurrenceRepository
	transactionRepo repository.TransactionRepository
	transactionService TransactionService
}

func NewRecurrenceService(
	recurrenceRepo repository.RecurrenceRepository,
	transactionRepo repository.TransactionRepository,
	transactionService TransactionService,
) RecurrenceService {
	return &recurrenceService{
		recurrenceRepo:     recurrenceRepo,
		transactionRepo:    transactionRepo,
		transactionService: transactionService,
	}
}

func (s *recurrenceService) CreatePattern(ctx context.Context, userID uuid.UUID, req *dto.CreateRecurrencePatternRequest) (*dto.RecurrencePatternDTO, error) {
	transactionID, err := uuid.Parse(req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction_id inválido: %w", err)
	}

	// Verificar se a transação existe e pertence ao usuário
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação: %w", err)
	}
	if transaction == nil {
		return nil, fmt.Errorf("transação %s: %w", req.TransactionID, appError.ErrTransactionNotFound)
	}
	if transaction.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Verificar se já existe um padrão para esta transação
	existing, err := s.recurrenceRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar padrão existente: %w", err)
	}
	if existing != nil {
		return nil, errors.New("já existe um padrão de recorrência para esta transação")
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date inválida: %w", err)
		}
		endDate = &parsed
	}

	now := time.Now()
	pattern := &model.RecurrencePattern{
		ID:                uuid.New(),
		UserID:            userID,
		TransactionID:     transactionID,
		Frequency:         req.Frequency,
		Interval:          req.Interval,
		EndDate:           endDate,
		LastGeneratedDate: nil,
		CreatedAt:         now,
	}

	if err := s.recurrenceRepo.Create(ctx, pattern); err != nil {
		return nil, fmt.Errorf("erro ao criar padrão de recorrência: %w", err)
	}

	return s.modelToDTO(pattern), nil
}

func (s *recurrenceService) GetPatternByID(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) (*dto.RecurrencePatternDTO, error) {
	pattern, err := s.recurrenceRepo.FindByID(ctx, patternID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar padrão: %w", err)
	}
	if pattern == nil {
		return nil, fmt.Errorf("padrão %s: %w", patternID, appError.ErrRecurrencePatternNotFound)
	}
	if pattern.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(pattern), nil
}

func (s *recurrenceService) GetPatternsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.RecurrencePatternDTO, error) {
	patterns, err := s.recurrenceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar padrões: %w", err)
	}

	dtos := make([]*dto.RecurrencePatternDTO, len(patterns))
	for i, pattern := range patterns {
		dtos[i] = s.modelToDTO(pattern)
	}

	return dtos, nil
}

func (s *recurrenceService) UpdatePattern(ctx context.Context, userID uuid.UUID, patternID uuid.UUID, req *dto.UpdateRecurrencePatternRequest) (*dto.RecurrencePatternDTO, error) {
	pattern, err := s.recurrenceRepo.FindByID(ctx, patternID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar padrão: %w", err)
	}
	if pattern == nil {
		return nil, fmt.Errorf("padrão %s: %w", patternID, appError.ErrRecurrencePatternNotFound)
	}
	if pattern.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date inválida: %w", err)
		}
		endDate = &parsed
	}

	pattern.Frequency = req.Frequency
	pattern.Interval = req.Interval
	pattern.EndDate = endDate

	if err := s.recurrenceRepo.Update(ctx, pattern); err != nil {
		return nil, fmt.Errorf("erro ao atualizar padrão: %w", err)
	}

	return s.modelToDTO(pattern), nil
}

func (s *recurrenceService) DeletePattern(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) error {
	pattern, err := s.recurrenceRepo.FindByID(ctx, patternID)
	if err != nil {
		return fmt.Errorf("erro ao buscar padrão: %w", err)
	}
	if pattern == nil {
		return fmt.Errorf("padrão %s: %w", patternID, appError.ErrRecurrencePatternNotFound)
	}
	if pattern.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.recurrenceRepo.Delete(ctx, patternID); err != nil {
		return fmt.Errorf("erro ao deletar padrão: %w", err)
	}

	return nil
}

func (s *recurrenceService) GenerateTransactions(ctx context.Context, userID uuid.UUID, untilDate time.Time) (int, error) {
	// Buscar todos os padrões ativos
	patterns, err := s.recurrenceRepo.FindActivePatterns(ctx, untilDate)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar padrões ativos: %w", err)
	}

	generatedCount := 0

	for _, pattern := range patterns {
		if pattern.UserID != userID {
			continue
		}

		// Buscar a transação base
		baseTransaction, err := s.transactionRepo.FindByID(ctx, pattern.TransactionID)
		if err != nil || baseTransaction == nil {
			continue
		}

		// Calcular próximas datas
		nextDate := s.calculateNextDate(pattern, baseTransaction.Date)
		generated := 0

		for nextDate.Before(untilDate) || nextDate.Equal(untilDate) {
			// Verificar se não passou do end_date
			if pattern.EndDate != nil && nextDate.After(*pattern.EndDate) {
				break
			}

			// Verificar se já existe transação nesta data
			// (simplificação: não verificar duplicatas por enquanto)

			// Criar nova transação
			createReq := &dto.CreateTransactionRequest{
				AccountID:   baseTransaction.AccountID.String(),
				CategoryID:  nil,
				Type:        baseTransaction.Type,
				Amount:      baseTransaction.Amount,
				Description: baseTransaction.Description,
				Date:        nextDate.Format("2006-01-02"),
			}

			if baseTransaction.CategoryID != nil {
				categoryIDStr := baseTransaction.CategoryID.String()
				createReq.CategoryID = &categoryIDStr
			}

			_, err := s.transactionService.CreateTransaction(ctx, userID, createReq)
			if err != nil {
				// Log erro mas continua com outros padrões
				continue
			}

			generated++
			nextDate = s.calculateNextDate(pattern, nextDate)
		}

		if generated > 0 {
			// Atualizar last_generated_date
			lastDate := s.calculateNextDate(pattern, baseTransaction.Date)
			for i := 1; i < generated; i++ {
				lastDate = s.calculateNextDate(pattern, lastDate)
			}
			s.recurrenceRepo.UpdateLastGeneratedDate(ctx, pattern.ID, lastDate)
		}

		generatedCount += generated
	}

	return generatedCount, nil
}

func (s *recurrenceService) calculateNextDate(pattern *model.RecurrencePattern, fromDate time.Time) time.Time {
	switch pattern.Frequency {
	case "daily":
		return fromDate.AddDate(0, 0, pattern.Interval)
	case "weekly":
		return fromDate.AddDate(0, 0, 7*pattern.Interval)
	case "monthly":
		return fromDate.AddDate(0, pattern.Interval, 0)
	case "yearly":
		return fromDate.AddDate(pattern.Interval, 0, 0)
	default:
		return fromDate
	}
}

func (s *recurrenceService) modelToDTO(pattern *model.RecurrencePattern) *dto.RecurrencePatternDTO {
	dto := &dto.RecurrencePatternDTO{
		ID:        pattern.ID.String(),
		UserID:    pattern.UserID.String(),
		TransactionID: pattern.TransactionID.String(),
		Frequency: pattern.Frequency,
		Interval:  pattern.Interval,
		CreatedAt: pattern.CreatedAt.Format(time.RFC3339),
	}

	if pattern.EndDate != nil {
		endDateStr := pattern.EndDate.Format("2006-01-02")
		dto.EndDate = &endDateStr
	}

	if pattern.LastGeneratedDate != nil {
		lastGenStr := pattern.LastGeneratedDate.Format("2006-01-02")
		dto.LastGeneratedDate = &lastGenStr
	}

	return dto
}
