package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, userID uuid.UUID, req *dto.CreateTaskRequest) (*dto.TaskDTO, error)
	GetTaskByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*dto.TaskDTO, error)
	GetTasksByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TaskFilters) ([]*dto.TaskDTO, error)
	UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, req *dto.UpdateTaskRequest) (*dto.TaskDTO, error)
	DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error
	ReorderTask(ctx context.Context, userID uuid.UUID, req *dto.ReorderTasksRequest) (*dto.TaskDTO, error)
	CompleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*dto.TaskDTO, error)
	UncompleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*dto.TaskDTO, error)
}

type taskService struct {
	db              *sqlx.DB
	taskRepo        repository.TaskRepository
	statusRepo      repository.TaskStatusRepository
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewTaskService(
	db *sqlx.DB,
	taskRepo repository.TaskRepository,
	statusRepo repository.TaskStatusRepository,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
) TaskService {
	return &taskService{
		db:              db,
		taskRepo:        taskRepo,
		statusRepo:      statusRepo,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *taskService) CreateTask(ctx context.Context, userID uuid.UUID, req *dto.CreateTaskRequest) (*dto.TaskDTO, error) {
	statusID, err := uuid.Parse(req.StatusID)
	if err != nil {
		return nil, fmt.Errorf("status_id inválido: %w", err)
	}

	// Verificar se o status existe e pertence ao usuário
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

	// Gerar lexorank (no final da lista)
	maxLexorank, err := s.taskRepo.GetMaxLexorank(ctx, statusID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar max lexorank: %w", err)
	}
	newLexorank := GetNextLexorank(maxLexorank)

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("due_date inválida: %w", err)
		}
		dueDate = &parsed
	}

	var financialAccountID *uuid.UUID
	if req.FinancialAccountID != nil && *req.FinancialAccountID != "" {
		parsed, err := uuid.Parse(*req.FinancialAccountID)
		if err != nil {
			return nil, fmt.Errorf("financial_account_id inválido: %w", err)
		}
		financialAccountID = &parsed
	}

	var financialCategoryID *uuid.UUID
	if req.FinancialCategoryID != nil && *req.FinancialCategoryID != "" {
		parsed, err := uuid.Parse(*req.FinancialCategoryID)
		if err != nil {
			return nil, fmt.Errorf("financial_category_id inválido: %w", err)
		}
		financialCategoryID = &parsed
	}

	now := time.Now()
	task := &model.Task{
		ID:                  uuid.New(),
		UserID:              userID,
		StatusID:            statusID,
		Title:               req.Title,
		Description:         req.Description,
		Priority:            req.Priority,
		DueDate:             dueDate,
		Lexorank:            newLexorank,
		FinancialAccountID:  financialAccountID,
		FinancialAmount:     req.FinancialAmount,
		FinancialCategoryID: financialCategoryID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("erro ao criar tarefa: %w", err)
	}

	return s.modelToDTO(task), nil
}

func (s *taskService) GetTaskByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*dto.TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("tarefa %s: %w", taskID, appError.ErrTaskNotFound)
	}
	if task.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(task), nil
}

func (s *taskService) GetTasksByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TaskFilters) ([]*dto.TaskDTO, error) {
	tasks, err := s.taskRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
	}

	dtos := make([]*dto.TaskDTO, len(tasks))
	for i, task := range tasks {
		dtos[i] = s.modelToDTO(task)
	}

	return dtos, nil
}

func (s *taskService) UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, req *dto.UpdateTaskRequest) (*dto.TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("tarefa %s: %w", taskID, appError.ErrTaskNotFound)
	}
	if task.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	statusID, err := uuid.Parse(req.StatusID)
	if err != nil {
		return nil, fmt.Errorf("status_id inválido: %w", err)
	}

	// Verificar se o status existe e pertence ao usuário
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

	// Se mudou de status, atualizar lexorank para o final do novo status
	if task.StatusID != statusID {
		maxLexorank, err := s.taskRepo.GetMaxLexorank(ctx, statusID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar max lexorank: %w", err)
		}
		task.Lexorank = GetNextLexorank(maxLexorank)
		task.StatusID = statusID
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("due_date inválida: %w", err)
		}
		dueDate = &parsed
	}

	var financialAccountID *uuid.UUID
	if req.FinancialAccountID != nil && *req.FinancialAccountID != "" {
		parsed, err := uuid.Parse(*req.FinancialAccountID)
		if err != nil {
			return nil, fmt.Errorf("financial_account_id inválido: %w", err)
		}
		financialAccountID = &parsed
	}

	var financialCategoryID *uuid.UUID
	if req.FinancialCategoryID != nil && *req.FinancialCategoryID != "" {
		parsed, err := uuid.Parse(*req.FinancialCategoryID)
		if err != nil {
			return nil, fmt.Errorf("financial_category_id inválido: %w", err)
		}
		financialCategoryID = &parsed
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Priority = req.Priority
	task.DueDate = dueDate
	task.FinancialAccountID = financialAccountID
	task.FinancialAmount = req.FinancialAmount
	task.FinancialCategoryID = financialCategoryID
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("erro ao atualizar tarefa: %w", err)
	}

	return s.modelToDTO(task), nil
}

func (s *taskService) DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	if task == nil {
		return errors.New("tarefa não encontrada")
	}
	if task.UserID != userID {
		return errors.New("tarefa não pertence ao usuário")
	}

	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("erro ao deletar tarefa: %w", err)
	}

	return nil
}

func (s *taskService) ReorderTask(ctx context.Context, userID uuid.UUID, req *dto.ReorderTasksRequest) (*dto.TaskDTO, error) {
	taskID, err := uuid.Parse(req.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task_id inválido: %w", err)
	}

	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("tarefa %s: %w", taskID, appError.ErrTaskNotFound)
	}
	if task.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	statusID, err := uuid.Parse(req.StatusID)
	if err != nil {
		return nil, fmt.Errorf("status_id inválido: %w", err)
	}

	// Verificar se o status existe e pertence ao usuário
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

	// Calcular novo lexorank
	var newLexorank string
	if req.AfterID != nil && *req.AfterID != "" {
		afterID, err := uuid.Parse(*req.AfterID)
		if err != nil {
			return nil, fmt.Errorf("after_id inválido: %w", err)
		}

		afterTask, err := s.taskRepo.GetTaskAfter(ctx, statusID, afterID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar tarefa após: %w", err)
		}
		if afterTask == nil {
			return nil, fmt.Errorf("tarefa após %s: %w", afterID, appError.ErrTaskNotFound)
		}

		// Buscar a próxima tarefa após afterTask
		tasks, err := s.taskRepo.FindByStatusID(ctx, statusID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
		}

		var beforeLexorank string
		for i, t := range tasks {
			if t.ID == afterID {
				if i+1 < len(tasks) {
					beforeLexorank = tasks[i+1].Lexorank
				}
				break
			}
		}

		newLexorank = GenerateLexorank(afterTask.Lexorank, beforeLexorank)
	} else {
		// Inserir no início
		tasks, err := s.taskRepo.FindByStatusID(ctx, statusID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar tarefas: %w", err)
		}

		var beforeLexorank string
		if len(tasks) > 0 {
			beforeLexorank = tasks[0].Lexorank
		}

		newLexorank = GenerateLexorank("", beforeLexorank)
	}

	// Atualizar status e lexorank
	if task.StatusID != statusID {
		if err := s.taskRepo.UpdateStatus(ctx, taskID, statusID); err != nil {
			return nil, fmt.Errorf("erro ao atualizar status: %w", err)
		}
		task.StatusID = statusID
	}

	if err := s.taskRepo.UpdateLexorank(ctx, taskID, newLexorank); err != nil {
		return nil, fmt.Errorf("erro ao atualizar lexorank: %w", err)
	}
	task.Lexorank = newLexorank

	return s.modelToDTO(task), nil
}

func (s *taskService) CompleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*dto.TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("tarefa %s: %w", taskID, appError.ErrTaskNotFound)
	}
	if task.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Se a tarefa tem informações financeiras, criar transação
	if task.FinancialAccountID != nil && task.FinancialAmount != nil && *task.FinancialAmount > 0 {
		// Verificar se a conta existe e pertence ao usuário
		account, err := s.accountRepo.FindByID(ctx, *task.FinancialAccountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar conta: %w", err)
		}
		if account == nil {
			return nil, fmt.Errorf("conta %s: %w", *task.FinancialAccountID, appError.ErrAccountNotFound)
		}
		if account.UserID != userID {
			return nil, appError.ErrUnauthorizedAccess
		}

		// Usar transação ACID para garantir consistência
		tx, err := s.db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
		}
		defer tx.Rollback()

		// Criar transação de despesa (expense) quando a tarefa é completada
		transaction := &model.Transaction{
			ID:          uuid.New(),
			UserID:      userID,
			AccountID:   *task.FinancialAccountID,
			CategoryID:  task.FinancialCategoryID,
			Type:        "expense",
			Amount:      *task.FinancialAmount,
			Description: fmt.Sprintf("Tarefa: %s", task.Title),
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		}

		// Criar transação
		createQuery := `
			INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date,
			                         status, created_at)
			VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
			        'paid', :created_at)
		`
		_, err = tx.NamedExecContext(ctx, createQuery, transaction)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar transação: %w", err)
		}

		// Atualizar saldo da conta
		balanceAdjustment := -transaction.Amount
		updateBalanceQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateBalanceQuery, balanceAdjustment, transaction.AccountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
		}

		// Completar tarefa dentro da mesma transação
		completeQuery := `UPDATE tasks SET completed_at = $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, completeQuery, time.Now(), taskID)
		if err != nil {
			return nil, fmt.Errorf("erro ao completar tarefa: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("erro ao confirmar transação: %w", err)
		}

		task.CompletedAt = func() *time.Time { t := time.Now(); return &t }()
		return s.modelToDTO(task), nil
	}

	// Se não tem informações financeiras, apenas completar tarefa
	if err := s.taskRepo.Complete(ctx, taskID); err != nil {
		return nil, fmt.Errorf("erro ao completar tarefa: %w", err)
	}

	task.CompletedAt = func() *time.Time { t := time.Now(); return &t }()
	return s.modelToDTO(task), nil
}

func (s *taskService) UncompleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*dto.TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tarefa: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("tarefa %s: %w", taskID, appError.ErrTaskNotFound)
	}
	if task.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	if err := s.taskRepo.Uncomplete(ctx, taskID); err != nil {
		return nil, fmt.Errorf("erro ao descompletar tarefa: %w", err)
	}

	task.CompletedAt = nil
	return s.modelToDTO(task), nil
}

func (s *taskService) modelToDTO(task *model.Task) *dto.TaskDTO {
	dto := &dto.TaskDTO{
		ID:          task.ID.String(),
		UserID:      task.UserID.String(),
		StatusID:    task.StatusID.String(),
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Lexorank:    task.Lexorank,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}

	if task.DueDate != nil {
		dueDateStr := task.DueDate.Format("2006-01-02")
		dto.DueDate = &dueDateStr
	}

	if task.CompletedAt != nil {
		completedAtStr := task.CompletedAt.Format(time.RFC3339)
		dto.CompletedAt = &completedAtStr
	}

	if task.FinancialAccountID != nil {
		accountIDStr := task.FinancialAccountID.String()
		dto.FinancialAccountID = &accountIDStr
	}

	if task.FinancialAmount != nil {
		dto.FinancialAmount = task.FinancialAmount
	}

	if task.FinancialCategoryID != nil {
		categoryIDStr := task.FinancialCategoryID.String()
		dto.FinancialCategoryID = &categoryIDStr
	}

	return dto
}
