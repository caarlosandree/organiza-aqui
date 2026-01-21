package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type RecurrenceRepository interface {
	Create(ctx context.Context, pattern *model.RecurrencePattern) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.RecurrencePattern, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.RecurrencePattern, error)
	FindByTransactionID(ctx context.Context, transactionID uuid.UUID) (*model.RecurrencePattern, error)
	FindActivePatterns(ctx context.Context, untilDate time.Time) ([]*model.RecurrencePattern, error)
	Update(ctx context.Context, pattern *model.RecurrencePattern) error
	UpdateLastGeneratedDate(ctx context.Context, id uuid.UUID, date time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type recurrenceRepository struct {
	db *sqlx.DB
}

func NewRecurrenceRepository(db *sqlx.DB) RecurrenceRepository {
	return &recurrenceRepository{db: db}
}

func (r *recurrenceRepository) Create(ctx context.Context, pattern *model.RecurrencePattern) error {
	query := `
		INSERT INTO recurrence_patterns (id, user_id, transaction_id, frequency, interval, end_date, last_generated_date, created_at)
		VALUES (:id, :user_id, :transaction_id, :frequency, :interval, :end_date, :last_generated_date, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, pattern)
	return err
}

func (r *recurrenceRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.RecurrencePattern, error) {
	var pattern model.RecurrencePattern
	query := `
		SELECT id, user_id, transaction_id, frequency, interval, end_date, last_generated_date, created_at
		FROM recurrence_patterns
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &pattern, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pattern, nil
}

func (r *recurrenceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.RecurrencePattern, error) {
	var patterns []*model.RecurrencePattern
	query := `
		SELECT id, user_id, transaction_id, frequency, interval, end_date, last_generated_date, created_at
		FROM recurrence_patterns
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &patterns, query, userID)
	return patterns, err
}

func (r *recurrenceRepository) FindByTransactionID(ctx context.Context, transactionID uuid.UUID) (*model.RecurrencePattern, error) {
	var pattern model.RecurrencePattern
	query := `
		SELECT id, user_id, transaction_id, frequency, interval, end_date, last_generated_date, created_at
		FROM recurrence_patterns
		WHERE transaction_id = $1
	`
	err := r.db.GetContext(ctx, &pattern, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pattern, nil
}

func (r *recurrenceRepository) FindActivePatterns(ctx context.Context, untilDate time.Time) ([]*model.RecurrencePattern, error) {
	var patterns []*model.RecurrencePattern
	query := `
		SELECT id, user_id, transaction_id, frequency, interval, end_date, last_generated_date, created_at
		FROM recurrence_patterns
		WHERE (end_date IS NULL OR end_date >= $1)
		ORDER BY last_generated_date NULLS FIRST, created_at
	`
	err := r.db.SelectContext(ctx, &patterns, query, untilDate)
	return patterns, err
}

func (r *recurrenceRepository) Update(ctx context.Context, pattern *model.RecurrencePattern) error {
	query := `
		UPDATE recurrence_patterns
		SET frequency = :frequency, interval = :interval, end_date = :end_date
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, pattern)
	return err
}

func (r *recurrenceRepository) UpdateLastGeneratedDate(ctx context.Context, id uuid.UUID, date time.Time) error {
	query := `UPDATE recurrence_patterns SET last_generated_date = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, date, id)
	return err
}

func (r *recurrenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM recurrence_patterns WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
