package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type TaskStatusRepository interface {
	Create(ctx context.Context, status *model.TaskStatus) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.TaskStatus, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskStatus, error)
	FindDefaultByUserID(ctx context.Context, userID uuid.UUID) (*model.TaskStatus, error)
	Update(ctx context.Context, status *model.TaskStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateOrder(ctx context.Context, statuses []*model.TaskStatus) error
}

type taskStatusRepository struct {
	db *sqlx.DB
}

func NewTaskStatusRepository(db *sqlx.DB) TaskStatusRepository {
	return &taskStatusRepository{db: db}
}

func (r *taskStatusRepository) Create(ctx context.Context, status *model.TaskStatus) error {
	query := `
		INSERT INTO task_statuses (id, user_id, name, color, order_index, is_default, created_at, updated_at)
		VALUES (:id, :user_id, :name, :color, :order_index, :is_default, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, status)
	return err
}

func (r *taskStatusRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.TaskStatus, error) {
	var status model.TaskStatus
	query := `
		SELECT id, user_id, name, color, order_index, is_default, created_at, updated_at
		FROM task_statuses
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &status, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}

func (r *taskStatusRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskStatus, error) {
	var statuses []*model.TaskStatus
	query := `
		SELECT id, user_id, name, color, order_index, is_default, created_at, updated_at
		FROM task_statuses
		WHERE user_id = $1
		ORDER BY order_index ASC, created_at ASC
	`
	err := r.db.SelectContext(ctx, &statuses, query, userID)
	return statuses, err
}

func (r *taskStatusRepository) FindDefaultByUserID(ctx context.Context, userID uuid.UUID) (*model.TaskStatus, error) {
	var status model.TaskStatus
	query := `
		SELECT id, user_id, name, color, order_index, is_default, created_at, updated_at
		FROM task_statuses
		WHERE user_id = $1 AND is_default = true
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &status, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}

func (r *taskStatusRepository) Update(ctx context.Context, status *model.TaskStatus) error {
	query := `
		UPDATE task_statuses
		SET name = :name, color = :color, order_index = :order_index, is_default = :is_default, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, status)
	return err
}

func (r *taskStatusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM task_statuses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *taskStatusRepository) UpdateOrder(ctx context.Context, statuses []*model.TaskStatus) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, status := range statuses {
		query := `UPDATE task_statuses SET order_index = $1, updated_at = NOW() WHERE id = $2`
		if _, err := tx.ExecContext(ctx, query, status.OrderIndex, status.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
