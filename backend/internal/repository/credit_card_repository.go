package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type CreditCardRepository interface {
	Create(ctx context.Context, creditCard *model.CreditCard) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.CreditCard, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.CreditCard, error)
	Update(ctx context.Context, creditCard *model.CreditCard) error
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateUsedLimit(ctx context.Context, creditCardID uuid.UUID) (int64, error)
	FindOpenBill(ctx context.Context, creditCardID uuid.UUID) (*model.CreditCardBill, error)
}

type creditCardRepository struct {
	db *sqlx.DB
}

func NewCreditCardRepository(db *sqlx.DB) CreditCardRepository {
	return &creditCardRepository{db: db}
}

func (r *creditCardRepository) Create(ctx context.Context, creditCard *model.CreditCard) error {
	query := `
		INSERT INTO credit_cards (id, user_id, name, account_id, limit_amount, closing_day, due_day, color, created_at, updated_at)
		VALUES (:id, :user_id, :name, :account_id, :limit_amount, :closing_day, :due_day, :color, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, creditCard)
	return err
}

func (r *creditCardRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.CreditCard, error) {
	var creditCard model.CreditCard
	query := `SELECT id, user_id, name, account_id, limit_amount, closing_day, due_day, color, created_at, updated_at 
	          FROM credit_cards WHERE id = $1`
	err := r.db.GetContext(ctx, &creditCard, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &creditCard, nil
}

func (r *creditCardRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.CreditCard, error) {
	var creditCards []*model.CreditCard
	query := `SELECT id, user_id, name, account_id, limit_amount, closing_day, due_day, color, created_at, updated_at 
	          FROM credit_cards WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &creditCards, query, userID)
	return creditCards, err
}

func (r *creditCardRepository) Update(ctx context.Context, creditCard *model.CreditCard) error {
	query := `
		UPDATE credit_cards 
		SET name = :name, account_id = :account_id, limit_amount = :limit_amount,
		    closing_day = :closing_day, due_day = :due_day, color = :color, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, creditCard)
	return err
}

func (r *creditCardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM credit_cards WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CalculateUsedLimit calcula o limite utilizado somando transações pendentes vinculadas ao cartão
func (r *creditCardRepository) CalculateUsedLimit(ctx context.Context, creditCardID uuid.UUID) (int64, error) {
	var usedLimit int64
	// Buscar a conta associada ao cartão
	var accountID uuid.UUID
	query := `SELECT account_id FROM credit_cards WHERE id = $1`
	err := r.db.GetContext(ctx, &accountID, query, creditCardID)
	if err != nil {
		return 0, err
	}

	// Calcular soma de transações pendentes na conta do cartão
	// Assumindo que transações na conta do cartão são gastos do cartão
	query2 := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE account_id = $1 
		  AND type = 'expense'
		  AND status = 'pending'
	`
	err = r.db.GetContext(ctx, &usedLimit, query2, accountID)
	return usedLimit, err
}

// FindOpenBill busca a fatura aberta (open) do cartão
func (r *creditCardRepository) FindOpenBill(ctx context.Context, creditCardID uuid.UUID) (*model.CreditCardBill, error) {
	var bill model.CreditCardBill
	query := `SELECT id, credit_card_id, month, year, status, closing_date, due_date, payment_transaction_id, created_at, updated_at
	          FROM credit_card_bills 
	          WHERE credit_card_id = $1 AND status = 'open'
	          ORDER BY year DESC, month DESC
	          LIMIT 1`
	err := r.db.GetContext(ctx, &bill, query, creditCardID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bill, nil
}
