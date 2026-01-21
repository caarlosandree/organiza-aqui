package repository

import (
	"context"
	"database/sql"
	"time"

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
// Pagamentos (income) reduzem o total da fatura
// O período da fatura vai do dia seguinte ao fechamento anterior até o fechamento atual
func (r *creditCardBillRepository) CalculateBillTotal(ctx context.Context, billID uuid.UUID) (int64, error) {
	// Buscar informações da fatura
	var bill model.CreditCardBill
	query := `SELECT credit_card_id, month, year, closing_date FROM credit_card_bills WHERE id = $1`
	err := r.db.GetContext(ctx, &bill, query, billID)
	if err != nil {
		return 0, err
	}

	// Buscar a conta do cartão e o dia de fechamento
	var accountID uuid.UUID
	var closingDay int
	query2 := `SELECT account_id, closing_day FROM credit_cards WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query2, bill.CreditCardID).Scan(&accountID, &closingDay)
	if err != nil {
		return 0, err
	}

	// Calcular data de início do período: dia seguinte ao fechamento anterior
	// Se a closing_date é 02/02/2026, o período começa em 03/01/2026
	// (dia seguinte ao fechamento do mês anterior, que seria 02/01/2026)
	// Calcular a data de fechamento do mês anterior
	prevMonth := bill.ClosingDate.AddDate(0, -1, 0)
	// Usar o mesmo dia de fechamento no mês anterior
	// Se o mês anterior não tem esse dia (ex: 31 em fevereiro), usar o último dia do mês
	prevMonthClosingDate := time.Date(prevMonth.Year(), prevMonth.Month(), closingDay, 0, 0, 0, 0, prevMonth.Location())
	// Verificar se a data foi ajustada (mês anterior não tinha esse dia)
	if prevMonthClosingDate.Month() != prevMonth.Month() {
		// Usar o último dia do mês anterior
		lastDayOfPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month()+1, 0, 0, 0, 0, 0, prevMonth.Location())
		prevMonthClosingDate = lastDayOfPrevMonth
	}
	
	// Data de início: dia seguinte ao fechamento anterior
	startDate := prevMonthClosingDate.AddDate(0, 0, 1)

	// Calcular total: soma de despesas MENOS pagamentos (income) na conta do cartão
	// dentro do período da fatura (do dia seguinte ao fechamento anterior até o fechamento atual)
	// Pagamentos reduzem o total da fatura
	var total int64
	query3 := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total
		FROM transactions
		WHERE account_id = $1
		  AND date >= $2
		  AND date <= $3
	`
	err = r.db.GetContext(ctx, &total, query3, accountID, startDate, bill.ClosingDate)
	return total, err
}

// CloseBill fecha uma fatura (muda status de 'open' para 'closed')
func (r *creditCardBillRepository) CloseBill(ctx context.Context, billID uuid.UUID) error {
	query := `UPDATE credit_card_bills SET status = 'closed', updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, billID)
	return err
}
