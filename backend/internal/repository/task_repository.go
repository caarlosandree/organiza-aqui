package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TaskFilters) ([]*model.Task, error)
	FindByStatusID(ctx context.Context, statusID uuid.UUID) ([]*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	UpdateStatus(ctx context.Context, taskID uuid.UUID, statusID uuid.UUID) error
	UpdateLexorank(ctx context.Context, taskID uuid.UUID, lexorank string) error
	Complete(ctx context.Context, taskID uuid.UUID) error
	Uncomplete(ctx context.Context, taskID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMaxLexorank(ctx context.Context, statusID uuid.UUID) (string, error)
	GetTaskAfter(ctx context.Context, statusID uuid.UUID, afterID uuid.UUID) (*model.Task, error)
}

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (id, user_id, status_id, title, description, priority, due_date, lexorank, 
		                  financial_account_id, financial_amount, financial_category_id, created_at, updated_at)
		VALUES (:id, :user_id, :status_id, :title, :description, :priority, :due_date, :lexorank, 
		        :financial_account_id, :financial_amount, :financial_category_id, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *taskRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task
	query := `
		SELECT id, user_id, status_id, title, description, priority, due_date, completed_at, lexorank, 
		       financial_account_id, financial_amount, financial_category_id, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TaskFilters) ([]*model.Task, error) {
	qb := NewQueryBuilder(`SELECT id, user_id, status_id, title, description, priority, due_date, completed_at, lexorank, 
	         financial_account_id, financial_amount, financial_category_id, created_at, updated_at 
	         FROM tasks`).
		WhereEqual("user_id", userID)

	if filters != nil {
		qb.WhereEqual("status_id", filters.StatusID).
			WhereEqual("priority", filters.Priority).
			WhereGreaterOrEqual("due_date", filters.StartDate).
			WhereLessOrEqual("due_date", filters.EndDate)

		if filters.Completed != nil {
			if *filters.Completed {
				qb.WhereIsNotNull("completed_at")
			} else {
				qb.WhereIsNull("completed_at")
			}
		}

		qb.OrderBy("status_id, lexorank ASC, created_at ASC").
			Limit(filters.Limit).
			Offset(filters.Offset)
	} else {
		qb.OrderBy("status_id, lexorank ASC, created_at ASC")
	}

	query, args := qb.Build()
	var tasks []*model.Task
	err := r.db.SelectContext(ctx, &tasks, query, args...)
	return tasks, err
}

func (r *taskRepository) FindByStatusID(ctx context.Context, statusID uuid.UUID) ([]*model.Task, error) {
	var tasks []*model.Task
	query := `
		SELECT id, user_id, status_id, title, description, priority, due_date, completed_at, lexorank, 
		       financial_account_id, financial_amount, financial_category_id, created_at, updated_at
		FROM tasks
		WHERE status_id = $1
		ORDER BY lexorank ASC, created_at ASC
	`
	err := r.db.SelectContext(ctx, &tasks, query, statusID)
	return tasks, err
}

func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE tasks
		SET status_id = :status_id, title = :title, description = :description, priority = :priority, 
		    due_date = :due_date, financial_account_id = :financial_account_id, 
		    financial_amount = :financial_amount, financial_category_id = :financial_category_id,
		    updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *taskRepository) UpdateStatus(ctx context.Context, taskID uuid.UUID, statusID uuid.UUID) error {
	query := `UPDATE tasks SET status_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, statusID, taskID)
	return err
}

func (r *taskRepository) UpdateLexorank(ctx context.Context, taskID uuid.UUID, lexorank string) error {
	query := `UPDATE tasks SET lexorank = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, lexorank, taskID)
	return err
}

func (r *taskRepository) Complete(ctx context.Context, taskID uuid.UUID) error {
	query := `UPDATE tasks SET completed_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, taskID)
	return err
}

func (r *taskRepository) Uncomplete(ctx context.Context, taskID uuid.UUID) error {
	query := `UPDATE tasks SET completed_at = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, taskID)
	return err
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *taskRepository) GetMaxLexorank(ctx context.Context, statusID uuid.UUID) (string, error) {
	var maxLexorank string
	query := `SELECT COALESCE(MAX(lexorank), '0|') FROM tasks WHERE status_id = $1`
	err := r.db.GetContext(ctx, &maxLexorank, query, statusID)
	return maxLexorank, err
}

func (r *taskRepository) GetTaskAfter(ctx context.Context, statusID uuid.UUID, afterID uuid.UUID) (*model.Task, error) {
	var task model.Task
	query := `
		SELECT id, user_id, status_id, title, description, priority, due_date, completed_at, lexorank, 
		       financial_account_id, financial_amount, financial_category_id, created_at, updated_at
		FROM tasks
		WHERE status_id = $1 AND id = $2
	`
	err := r.db.GetContext(ctx, &task, query, statusID, afterID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}
