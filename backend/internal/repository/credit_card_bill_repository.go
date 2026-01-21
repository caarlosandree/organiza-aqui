package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type CreditCardBillRepository interface {
	Create(ctx context.Context, bill *model.CreditCardBill) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.CreditCardBill, error)
	FindByCreditCardID(ctx context.Context, creditCardID uuid.UUID) ([]*model.CreditCardBill, error)
	FindByPeriod(ctx context.Context, creditCardID uuid.UUID, year int, month int) (*model.CreditCardBill, error)
	Update(ctx context.Context, bill *model.CreditCardBill) error
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateBillTotal(ctx context.Context, billID uuid.UUID) (int64, error)
	CloseBill(ctx context.Context, billID uuid.UUID) error
}

type creditCardBillRepository struct {
	db *sqlx.DB
}

func NewCreditCardBillRepository(db *sqlx.DB) CreditCardBillRepository {
	return &creditCardBillRepository{db: db}
}

func (r *creditCardBillRepository) Create(ctx context.Context, bill *model.CreditCardBill) error {
	query := `
		INSERT INTO credit_card_bills (id, credit_card_id, month, year, status, closing_date, due_date, payment_transaction_id, created_at, updated_at)
		VALUES (:id, :credit_card_id, :month, :year, :status, :closing_date, :due_date, :payment_transaction_id, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, bill)
	return err
}

func (r *creditCardBillRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.CreditCardBill, error) {
	var bill model.CreditCardBill
	query := `SELECT id, credit_card_id, month, year, status, closing_date, due_date, payment_transaction_id, created_at, updated_at
	          FROM credit_card_bills WHERE id = $1`
	err := r.db.GetContext(ctx, &bill, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bill, nil
}

func (r *creditCardBillRepository) FindByCreditCardID(ctx context.Context, creditCardID uuid.UUID) ([]*model.CreditCardBill, error) {
	var bills []*model.CreditCardBill
	query := `SELECT id, credit_card_id, month, year, status, closing_date, due_date, payment_transaction_id, created_at, updated_at
	          FROM credit_card_bills 
	          WHERE credit_card_id = $1 
	          ORDER BY year DESC, month DESC`
	err := r.db.SelectContext(ctx, &bills, query, creditCardID)
	return bills, err
}

func (r *creditCardBillRepository) FindByPeriod(ctx context.Context, creditCardID uuid.UUID, year int, month int) (*model.CreditCardBill, error) {
	var bill model.CreditCardBill
	query := `SELECT id, credit_card_id, month, year, status, closing_date, due_date, payment_transaction_id, created_at, updated_at
	          FROM credit_card_bills 
	          WHERE credit_card_id = $1 AND year = $2 AND month = $3`
	err := r.db.GetContext(ctx, &bill, query, creditCardID, year, month)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bill, nil
}

func (r *creditCardBillRepository) Update(ctx context.Context, bill *model.CreditCardBill) error {
	query := `
		UPDATE credit_card_bills 
		SET status = :status, payment_transaction_id = :payment_transaction_id, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, bill)
	return err
}

func (r *creditCardBillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM credit_card_bills WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CalculateBillTotal calcula o total da fatura somando transações vinculadas
// Assumimos que transações na conta do cartão dentro do período da fatura são parte da fatura
func (r *creditCardBillRepository) CalculateBillTotal(ctx context.Context, billID uuid.UUID) (int64, error) {
	// Buscar informações da fatura
	var bill model.CreditCardBill
	query := `SELECT credit_card_id, month, year, closing_date FROM credit_card_bills WHERE id = $1`
	err := r.db.GetContext(ctx, &bill, query, billID)
	if err != nil {
		return 0, err
	}

	// Buscar a conta do cartão
	var accountID uuid.UUID
	query2 := `SELECT account_id FROM credit_cards WHERE id = $1`
	err = r.db.GetContext(ctx, &accountID, query2, bill.CreditCardID)
	if err != nil {
		return 0, err
	}

	// Calcular total: soma de transações de despesa na conta do cartão
	// dentro do período da fatura (até a data de fechamento)
	var total int64
	query3 := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE account_id = $1
		  AND type = 'expense'
		  AND date <= $2
		  AND date >= DATE_TRUNC('month', $2::date)
	`
	err = r.db.GetContext(ctx, &total, query3, accountID, bill.ClosingDate)
	return total, err
}

// CloseBill fecha uma fatura (muda status de 'open' para 'closed')
func (r *creditCardBillRepository) CloseBill(ctx context.Context, billID uuid.UUID) error {
	query := `UPDATE credit_card_bills SET status = 'closed', updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, billID)
	return err
}
