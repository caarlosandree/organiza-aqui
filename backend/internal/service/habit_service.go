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

type HabitService interface {
	CreateHabit(ctx context.Context, userID uuid.UUID, req *dto.CreateHabitRequest) (*dto.HabitDTO, error)
	GetHabitByID(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) (*dto.HabitDTO, error)
	GetHabitsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.HabitDTO, error)
	UpdateHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID, req *dto.UpdateHabitRequest) (*dto.HabitDTO, error)
	DeleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) error
	GetHabitStats(ctx context.Context, userID uuid.UUID, habitID uuid.UUID, startDate, endDate time.Time) (*dto.HabitStatsDTO, error)
}

type HabitTrackingService interface {
	CreateTracking(ctx context.Context, userID uuid.UUID, req *dto.CreateHabitTrackingRequest) (*dto.HabitTrackingDTO, error)
	GetTrackingByID(ctx context.Context, userID uuid.UUID, trackingID uuid.UUID) (*dto.HabitTrackingDTO, error)
	GetTrackingByHabitID(ctx context.Context, userID uuid.UUID, habitID uuid.UUID, startDate, endDate time.Time) ([]*dto.HabitTrackingDTO, error)
	UpdateTracking(ctx context.Context, userID uuid.UUID, trackingID uuid.UUID, req *dto.UpdateHabitTrackingRequest) (*dto.HabitTrackingDTO, error)
	DeleteTracking(ctx context.Context, userID uuid.UUID, trackingID uuid.UUID) error
}

type habitService struct {
	habitRepo  repository.HabitRepository
	trackingRepo repository.HabitTrackingRepository
}

func NewHabitService(habitRepo repository.HabitRepository, trackingRepo repository.HabitTrackingRepository) HabitService {
	return &habitService{
		habitRepo:    habitRepo,
		trackingRepo: trackingRepo,
	}
}

func (s *habitService) CreateHabit(ctx context.Context, userID uuid.UUID, req *dto.CreateHabitRequest) (*dto.HabitDTO, error) {
	now := time.Now()
	habit := &model.Habit{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Frequency:   req.Frequency,
		TargetDays:  req.TargetDays,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.habitRepo.Create(ctx, habit); err != nil {
		return nil, fmt.Errorf("erro ao criar hábito: %w", err)
	}

	return s.modelToDTO(habit), nil
}

func (s *habitService) GetHabitByID(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) (*dto.HabitDTO, error) {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("hábito %s: %w", habitID, appError.ErrHabitNotFound)
	}
	if habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(habit), nil
}

func (s *habitService) GetHabitsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.HabitDTO, error) {
	habits, err := s.habitRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábitos: %w", err)
	}

	dtos := make([]*dto.HabitDTO, len(habits))
	for i, habit := range habits {
		dtos[i] = s.modelToDTO(habit)
	}

	return dtos, nil
}

func (s *habitService) UpdateHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID, req *dto.UpdateHabitRequest) (*dto.HabitDTO, error) {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("hábito %s: %w", habitID, appError.ErrHabitNotFound)
	}
	if habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	habit.Name = req.Name
	habit.Description = req.Description
	habit.Color = req.Color
	habit.Frequency = req.Frequency
	habit.TargetDays = req.TargetDays
	habit.UpdatedAt = time.Now()

	if err := s.habitRepo.Update(ctx, habit); err != nil {
		return nil, fmt.Errorf("erro ao atualizar hábito: %w", err)
	}

	return s.modelToDTO(habit), nil
}

func (s *habitService) DeleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) error {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil {
		return fmt.Errorf("hábito %s: %w", habitID, appError.ErrHabitNotFound)
	}
	if habit.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.habitRepo.Delete(ctx, habitID); err != nil {
		return fmt.Errorf("erro ao deletar hábito: %w", err)
	}

	return nil
}

func (s *habitService) GetHabitStats(ctx context.Context, userID uuid.UUID, habitID uuid.UUID, startDate, endDate time.Time) (*dto.HabitStatsDTO, error) {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("hábito %s: %w", habitID, appError.ErrHabitNotFound)
	}
	if habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	stats, err := s.trackingRepo.GetStats(ctx, habitID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %w", err)
	}

	return &dto.HabitStatsDTO{
		HabitID:        habitID.String(),
		TotalDays:      stats.TotalDays,
		CompletedDays:  stats.CompletedDays,
		CompletionRate: stats.CompletionRate,
		CurrentStreak:  stats.CurrentStreak,
		LongestStreak:  stats.LongestStreak,
	}, nil
}

func (s *habitService) modelToDTO(habit *model.Habit) *dto.HabitDTO {
	return &dto.HabitDTO{
		ID:          habit.ID.String(),
		UserID:      habit.UserID.String(),
		Name:        habit.Name,
		Description: habit.Description,
		Color:       habit.Color,
		Frequency:   habit.Frequency,
		TargetDays:  habit.TargetDays,
		CreatedAt:   habit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   habit.UpdatedAt.Format(time.RFC3339),
	}
}

type habitTrackingService struct {
	trackingRepo repository.HabitTrackingRepository
	habitRepo    repository.HabitRepository
}

func NewHabitTrackingService(trackingRepo repository.HabitTrackingRepository, habitRepo repository.HabitRepository) HabitTrackingService {
	return &habitTrackingService{
		trackingRepo: trackingRepo,
		habitRepo:    habitRepo,
	}
}

func (s *habitTrackingService) CreateTracking(ctx context.Context, userID uuid.UUID, req *dto.CreateHabitTrackingRequest) (*dto.HabitTrackingDTO, error) {
	habitID, err := uuid.Parse(req.HabitID)
	if err != nil {
		return nil, fmt.Errorf("habit_id inválido: %w", err)
	}

	// Verificar se o hábito existe e pertence ao usuário
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("hábito %s: %w", habitID, appError.ErrHabitNotFound)
	}
	if habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("date inválida: %w", err)
	}

	now := time.Now()
	tracking := &model.HabitTracking{
		ID:        uuid.New(),
		HabitID:   habitID,
		Date:      date,
		Completed: req.Completed,
		Notes:     req.Notes,
		CreatedAt: now,
	}

	if err := s.trackingRepo.Create(ctx, tracking); err != nil {
		return nil, fmt.Errorf("erro ao criar tracking: %w", err)
	}

	return s.modelToDTO(tracking), nil
}

func (s *habitTrackingService) GetTrackingByID(ctx context.Context, userID uuid.UUID, trackingID uuid.UUID) (*dto.HabitTrackingDTO, error) {
	tracking, err := s.trackingRepo.FindByID(ctx, trackingID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tracking: %w", err)
	}
	if tracking == nil {
		return nil, fmt.Errorf("tracking não encontrado: %w", appError.ErrTrackingNotFound)
	}

	// Verificar se o hábito pertence ao usuário
	habit, err := s.habitRepo.FindByID(ctx, tracking.HabitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil || habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(tracking), nil
}

func (s *habitTrackingService) GetTrackingByHabitID(ctx context.Context, userID uuid.UUID, habitID uuid.UUID, startDate, endDate time.Time) ([]*dto.HabitTrackingDTO, error) {
	// Verificar se o hábito pertence ao usuário
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("hábito %s: %w", habitID, appError.ErrHabitNotFound)
	}
	if habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	trackings, err := s.trackingRepo.FindByHabitID(ctx, habitID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar trackings: %w", err)
	}

	dtos := make([]*dto.HabitTrackingDTO, len(trackings))
	for i, tracking := range trackings {
		dtos[i] = s.modelToDTO(tracking)
	}

	return dtos, nil
}

func (s *habitTrackingService) UpdateTracking(ctx context.Context, userID uuid.UUID, trackingID uuid.UUID, req *dto.UpdateHabitTrackingRequest) (*dto.HabitTrackingDTO, error) {
	tracking, err := s.trackingRepo.FindByID(ctx, trackingID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tracking: %w", err)
	}
	if tracking == nil {
		return nil, fmt.Errorf("tracking não encontrado: %w", appError.ErrTrackingNotFound)
	}

	// Verificar se o hábito pertence ao usuário
	habit, err := s.habitRepo.FindByID(ctx, tracking.HabitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil || habit.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	tracking.Completed = req.Completed
	tracking.Notes = req.Notes

	if err := s.trackingRepo.Update(ctx, tracking); err != nil {
		return nil, fmt.Errorf("erro ao atualizar tracking: %w", err)
	}

	return s.modelToDTO(tracking), nil
}

func (s *habitTrackingService) DeleteTracking(ctx context.Context, userID uuid.UUID, trackingID uuid.UUID) error {
	tracking, err := s.trackingRepo.FindByID(ctx, trackingID)
	if err != nil {
		return fmt.Errorf("erro ao buscar tracking: %w", err)
	}
	if tracking == nil {
		return fmt.Errorf("tracking não encontrado: %w", appError.ErrTrackingNotFound)
	}

	// Verificar se o hábito pertence ao usuário
	habit, err := s.habitRepo.FindByID(ctx, tracking.HabitID)
	if err != nil {
		return fmt.Errorf("erro ao buscar hábito: %w", err)
	}
	if habit == nil || habit.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.trackingRepo.Delete(ctx, trackingID); err != nil {
		return fmt.Errorf("erro ao deletar tracking: %w", err)
	}

	return nil
}

func (s *habitTrackingService) modelToDTO(tracking *model.HabitTracking) *dto.HabitTrackingDTO {
	dto := &dto.HabitTrackingDTO{
		ID:        tracking.ID.String(),
		HabitID:   tracking.HabitID.String(),
		Date:      tracking.Date.Format("2006-01-02"),
		Completed: tracking.Completed,
		CreatedAt: tracking.CreatedAt.Format(time.RFC3339),
	}

	if tracking.Notes != "" {
		dto.Notes = tracking.Notes
	}

	return dto
}
