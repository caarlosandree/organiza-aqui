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

type TaskStatusService interface {
	CreateStatus(ctx context.Context, userID uuid.UUID, req *dto.CreateTaskStatusRequest) (*dto.TaskStatusDTO, error)
	GetStatusByID(ctx context.Context, userID uuid.UUID, statusID uuid.UUID) (*dto.TaskStatusDTO, error)
	GetStatusesByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.TaskStatusDTO, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, statusID uuid.UUID, req *dto.UpdateTaskStatusRequest) (*dto.TaskStatusDTO, error)
	DeleteStatus(ctx context.Context, userID uuid.UUID, statusID uuid.UUID) error
	ReorderStatuses(ctx context.Context, userID uuid.UUID, statusIDs []string) error
}

type taskStatusService struct {
	statusRepo repository.TaskStatusRepository
	taskRepo   repository.TaskRepository
}

func NewTaskStatusService(statusRepo repository.TaskStatusRepository, taskRepo repository.TaskRepository) TaskStatusService {
	return &taskStatusService{
		statusRepo: statusRepo,
		taskRepo:   taskRepo,
	}
}

func (s *taskStatusService) CreateStatus(ctx context.Context, userID uuid.UUID, req *dto.CreateTaskStatusRequest) (*dto.TaskStatusDTO, error) {
	// Se is_default, remover default de outros status
	if req.IsDefault {
		defaultStatus, _ := s.statusRepo.FindDefaultByUserID(ctx, userID)
		if defaultStatus != nil {
			defaultStatus.IsDefault = false
			defaultStatus.UpdatedAt = time.Now()
			if err := s.statusRepo.Update(ctx, defaultStatus); err != nil {
				return nil, fmt.Errorf("erro ao remover default anterior: %w", err)
			}
		}
	}

	now := time.Now()
	status := &model.TaskStatus{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       req.Name,
		Color:      req.Color,
		OrderIndex: req.OrderIndex,
		IsDefault:  req.IsDefault,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.statusRepo.Create(ctx, status); err != nil {
		return nil, fmt.Errorf("erro ao criar status: %w", err)
	}

	return s.modelToDTO(status), nil
}

func (s *taskStatusService) GetStatusByID(ctx context.Context, userID uuid.UUID, statusID uuid.UUID) (*dto.TaskStatusDTO, error) {
	status, err := s.statusRepo.FindByID(ctx, statusID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar status: %w", err)
	}
	if status == nil {
		return nil, fmt.Errorf("status %s: %w", statusID, appError.ErrStatusNotFound)
	}
	if status.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(status), nil
}

func (s *taskStatusService) GetStatusesByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.TaskStatusDTO, error) {
	statuses, err := s.statusRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar statuses: %w", err)
	}

	dtos := make([]*dto.TaskStatusDTO, len(statuses))
	for i, status := range statuses {
		dtos[i] = s.modelToDTO(status)
	}

	return dtos, nil
}

func (s *taskStatusService) UpdateStatus(ctx context.Context, userID uuid.UUID, statusID uuid.UUID, req *dto.UpdateTaskStatusRequest) (*dto.TaskStatusDTO, error) {
	status, err := s.statusRepo.FindByID(ctx, statusID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar status: %w", err)
	}
	if status == nil {
		return nil, fmt.Errorf("status %s: %w", statusID, appError.ErrStatusNotFound)
	}
	if status.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Se is_default mudou para true, remover default de outros
	if req.IsDefault && !status.IsDefault {
		defaultStatus, _ := s.statusRepo.FindDefaultByUserID(ctx, userID)
		if defaultStatus != nil && defaultStatus.ID != statusID {
			defaultStatus.IsDefault = false
			defaultStatus.UpdatedAt = time.Now()
			if err := s.statusRepo.Update(ctx, defaultStatus); err != nil {
				return nil, fmt.Errorf("erro ao remover default anterior: %w", err)
			}
		}
	}

	status.Name = req.Name
	status.Color = req.Color
	status.OrderIndex = req.OrderIndex
	status.IsDefault = req.IsDefault
	status.UpdatedAt = time.Now()

	if err := s.statusRepo.Update(ctx, status); err != nil {
		return nil, fmt.Errorf("erro ao atualizar status: %w", err)
	}

	return s.modelToDTO(status), nil
}

func (s *taskStatusService) DeleteStatus(ctx context.Context, userID uuid.UUID, statusID uuid.UUID) error {
	status, err := s.statusRepo.FindByID(ctx, statusID)
	if err != nil {
		return fmt.Errorf("erro ao buscar status: %w", err)
	}
	if status == nil {
		return fmt.Errorf("status %s: %w", statusID, appError.ErrStatusNotFound)
	}
	if status.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	// Verificar se há tarefas usando este status
	tasks, err := s.taskRepo.FindByStatusID(ctx, statusID)
	if err != nil {
		return fmt.Errorf("erro ao verificar tarefas: %w", err)
	}
	if len(tasks) > 0 {
		return errors.New("não é possível deletar status com tarefas associadas")
	}

	if err := s.statusRepo.Delete(ctx, statusID); err != nil {
		return fmt.Errorf("erro ao deletar status: %w", err)
	}

	return nil
}

func (s *taskStatusService) ReorderStatuses(ctx context.Context, userID uuid.UUID, statusIDs []string) error {
	statuses := make([]*model.TaskStatus, len(statusIDs))
	for i, idStr := range statusIDs {
		statusID, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("ID inválido: %s", idStr)
		}

		status, err := s.statusRepo.FindByID(ctx, statusID)
		if err != nil {
			return fmt.Errorf("erro ao buscar status: %w", err)
		}
		if status == nil {
			return fmt.Errorf("status não encontrado: %s", idStr)
		}
		if status.UserID != userID {
			return fmt.Errorf("status não pertence ao usuário: %s", idStr)
		}

		status.OrderIndex = i
		statuses[i] = status
	}

	return s.statusRepo.UpdateOrder(ctx, statuses)
}

func (s *taskStatusService) modelToDTO(status *model.TaskStatus) *dto.TaskStatusDTO {
	return &dto.TaskStatusDTO{
		ID:         status.ID.String(),
		UserID:     status.UserID.String(),
		Name:       status.Name,
		Color:      status.Color,
		OrderIndex: status.OrderIndex,
		IsDefault:  status.IsDefault,
		CreatedAt:  status.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  status.UpdatedAt.Format(time.RFC3339),
	}
}
