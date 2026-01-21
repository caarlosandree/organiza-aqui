package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Account, error)
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64) error
	SetBalance(ctx context.Context, accountID uuid.UUID, balance int64) error
	UpdateInitialBalance(ctx context.Context, accountID uuid.UUID, balance int64, date time.Time) error
}

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, name, type, balance, currency, bank_id, initial_balance, initial_balance_date, created_at)
		VALUES (:id, :user_id, :name, :type, :balance, :currency, :bank_id, :initial_balance, :initial_balance_date, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, account)
	return err
}

func (r *accountRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	var account model.Account
	query := `SELECT id, user_id, name, type, balance, currency, bank_id, initial_balance, initial_balance_date, created_at FROM accounts WHERE id = $1`
	err := r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Account, error) {
	var accounts []*model.Account
	query := `SELECT id, user_id, name, type, balance, currency, bank_id, initial_balance, initial_balance_date, created_at FROM accounts WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &accounts, query, userID)
	return accounts, err
}

func (r *accountRepository) Update(ctx context.Context, account *model.Account) error {
	query := `
		UPDATE accounts 
		SET name = :name, type = :type, currency = :currency, bank_id = :bank_id, initial_balance = :initial_balance, initial_balance_date = :initial_balance_date
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, account)
	return err
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *accountRepository) UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, amount, accountID)
	return err
}

func (r *accountRepository) SetBalance(ctx context.Context, accountID uuid.UUID, balance int64) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, balance, accountID)
	return err
}

func (r *accountRepository) UpdateInitialBalance(ctx context.Context, accountID uuid.UUID, balance int64, date time.Time) error {
	query := `UPDATE accounts SET initial_balance = $1, initial_balance_date = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, balance, date, accountID)
	return err
}
