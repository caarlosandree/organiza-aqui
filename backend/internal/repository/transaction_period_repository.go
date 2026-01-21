package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
)

type TransactionPeriodRepository interface {
	Create(ctx context.Context, period *model.TransactionPeriod) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.TransactionPeriod, error)
	FindByAccountAndPeriod(ctx context.Context, accountID uuid.UUID, periodType string, year, month int) (*model.TransactionPeriod, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionPeriodFilters) ([]*model.TransactionPeriod, error)
	Update(ctx context.Context, period *model.TransactionPeriod) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type transactionPeriodRepository struct {
	db *sqlx.DB
}

func NewTransactionPeriodRepository(db *sqlx.DB) TransactionPeriodRepository {
	return &transactionPeriodRepository{db: db}
}

func (r *transactionPeriodRepository) Create(ctx context.Context, period *model.TransactionPeriod) error {
	query := `
		INSERT INTO transaction_periods (id, user_id, account_id, period_type, year, month, status, created_at, updated_at)
		VALUES (:id, :user_id, :account_id, :period_type, :year, :month, :status, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, period)
	return err
}

func (r *transactionPeriodRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.TransactionPeriod, error) {
	var period model.TransactionPeriod
	query := `SELECT id, user_id, account_id, period_type, year, month, status, created_at, updated_at
	          FROM transaction_periods WHERE id = $1`
	err := r.db.GetContext(ctx, &period, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &period, nil
}

func (r *transactionPeriodRepository) FindByAccountAndPeriod(ctx context.Context, accountID uuid.UUID, periodType string, year, month int) (*model.TransactionPeriod, error) {
	var period model.TransactionPeriod
	query := `SELECT id, user_id, account_id, period_type, year, month, status, created_at, updated_at
	          FROM transaction_periods 
	          WHERE account_id = $1 AND period_type = $2 AND year = $3 AND month = $4`
	err := r.db.GetContext(ctx, &period, query, accountID, periodType, year, month)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &period, nil
}

func (r *transactionPeriodRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionPeriodFilters) ([]*model.TransactionPeriod, error) {
	qb := NewQueryBuilder(`
		SELECT tp.id, tp.user_id, tp.account_id, tp.period_type, tp.year, tp.month, tp.status, tp.created_at, tp.updated_at
		FROM transaction_periods tp
	`).
		WhereEqual("tp.user_id", userID)

	if filters != nil {
		qb.WhereEqual("tp.account_id", filters.AccountID).
			WhereEqual("tp.period_type", filters.PeriodType).
			WhereEqual("tp.year", filters.Year).
			WhereEqual("tp.month", filters.Month).
			WhereEqual("tp.status", filters.Status)
	}

	qb.OrderBy("tp.year DESC, tp.month DESC, tp.created_at DESC")

	query, args := qb.Build()
	var periods []*model.TransactionPeriod
	err := r.db.SelectContext(ctx, &periods, query, args...)
	return periods, err
}

func (r *transactionPeriodRepository) Update(ctx context.Context, period *model.TransactionPeriod) error {
	query := `
		UPDATE transaction_periods 
		SET status = :status, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, period)
	return err
}

func (r *transactionPeriodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM transaction_periods WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
