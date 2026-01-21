package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionFilters) ([]*model.Transaction, error)
	CountByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionFilters) (int, error)
	Update(ctx context.Context, transaction *model.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAccountBalanceAtDate(ctx context.Context, accountID uuid.UUID, date time.Time) (int64, error)
	GetAccountInitialBalance(ctx context.Context, accountID uuid.UUID) (int64, error)
	CalculateAccountBalanceFromTransactions(ctx context.Context, accountID uuid.UUID) (int64, error)
	FindByStatus(ctx context.Context, userID uuid.UUID, status string) ([]*model.Transaction, error)
	FindByTags(ctx context.Context, userID uuid.UUID, tags []string) ([]*model.Transaction, error)
	FindInstallments(ctx context.Context, parentTransactionID uuid.UUID) ([]*model.Transaction, error)
	FindByExternalID(ctx context.Context, externalID string) (*model.Transaction, error)
	CreateTransfer(ctx context.Context, fromTransaction, toTransaction *model.Transaction, fromAccountID, toAccountID uuid.UUID, amount int64) error
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date, 
		                         status, tags, to_account_id, parent_transaction_id, installment_number, 
		                         total_installments, external_id, period_id, reference_month, created_at)
		VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
		        :status, :tags, :to_account_id, :parent_transaction_id, :installment_number,
		        :total_installments, :external_id, :period_id, :reference_month, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, transaction)
	return err
}

func (r *transactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	var transaction model.Transaction
	query := `SELECT id, user_id, account_id, category_id, type, amount, description, date, 
	                 status, tags, to_account_id, parent_transaction_id, installment_number,
	                 total_installments, external_id, period_id, reference_month, created_at 
	          FROM transactions WHERE id = $1`
	err := r.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionFilters) ([]*model.Transaction, error) {
	qb := NewQueryBuilder(`
		SELECT id, user_id, account_id, category_id, type, amount, description, date,
		       status, tags, to_account_id, parent_transaction_id, installment_number,
		       total_installments, external_id, period_id, reference_month, created_at 
		FROM transactions
	`).
		WhereEqual("user_id", userID)

	if filters != nil {
		// Converter strings para UUID quando necessário
		if filters.AccountID != nil {
			accountID, err := uuid.Parse(*filters.AccountID)
			if err != nil {
				return nil, fmt.Errorf("account_id inválido: %w", err)
			}
			qb.WhereEqual("account_id", accountID)
		}

		if filters.CategoryID != nil {
			categoryID, err := uuid.Parse(*filters.CategoryID)
			if err != nil {
				return nil, fmt.Errorf("category_id inválido: %w", err)
			}
			qb.WhereEqual("category_id", categoryID)
		}

		if filters.ParentTransactionID != nil {
			parentTransactionID, err := uuid.Parse(*filters.ParentTransactionID)
			if err != nil {
				return nil, fmt.Errorf("parent_transaction_id inválido: %w", err)
			}
			qb.WhereEqual("parent_transaction_id", parentTransactionID)
		}

		if filters.Type != nil {
			qb.WhereEqual("type", *filters.Type)
		}

		if filters.Status != nil {
			qb.WhereEqual("status", *filters.Status)
		}

		if filters.StartDate != nil {
			qb.WhereGreaterOrEqual("date", *filters.StartDate)
		}

		if filters.EndDate != nil {
			qb.WhereLessOrEqual("date", *filters.EndDate)
		}

		if len(filters.Tags) > 0 {
			qb.WhereArrayOverlaps("tags", filters.Tags)
		}

		qb.OrderBy("date DESC, created_at DESC").
			Limit(filters.Limit).
			Offset(filters.Offset)
	} else {
		qb.OrderBy("date DESC, created_at DESC")
	}

	query, args := qb.Build()
	var transactions []*model.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, args...)
	return transactions, err
}

func (r *transactionRepository) CountByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionFilters) (int, error) {
	qb := NewQueryBuilder(`SELECT COUNT(*) FROM transactions`).
		WhereEqual("user_id", userID)

	if filters != nil {
		// Converter strings para UUID quando necessário
		if filters.AccountID != nil {
			accountID, err := uuid.Parse(*filters.AccountID)
			if err != nil {
				return 0, fmt.Errorf("account_id inválido: %w", err)
			}
			qb.WhereEqual("account_id", accountID)
		}

		if filters.CategoryID != nil {
			categoryID, err := uuid.Parse(*filters.CategoryID)
			if err != nil {
				return 0, fmt.Errorf("category_id inválido: %w", err)
			}
			qb.WhereEqual("category_id", categoryID)
		}

		if filters.ParentTransactionID != nil {
			parentTransactionID, err := uuid.Parse(*filters.ParentTransactionID)
			if err != nil {
				return 0, fmt.Errorf("parent_transaction_id inválido: %w", err)
			}
			qb.WhereEqual("parent_transaction_id", parentTransactionID)
		}

		if filters.Type != nil {
			qb.WhereEqual("type", *filters.Type)
		}

		if filters.Status != nil {
			qb.WhereEqual("status", *filters.Status)
		}

		if filters.StartDate != nil {
			qb.WhereGreaterOrEqual("date", *filters.StartDate)
		}

		if filters.EndDate != nil {
			qb.WhereLessOrEqual("date", *filters.EndDate)
		}

		if filters.MinAmount != nil {
			qb.WhereGreaterOrEqual("amount", *filters.MinAmount)
		}

		if filters.MaxAmount != nil {
			qb.WhereLessOrEqual("amount", *filters.MaxAmount)
		}

		if len(filters.Tags) > 0 {
			qb.WhereArrayOverlaps("tags", filters.Tags)
		}
	}

	query, args := qb.Build()
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

// GetAccountBalanceAtDate retorna o saldo de uma conta em uma data específica
func (r *transactionRepository) GetAccountBalanceAtDate(ctx context.Context, accountID uuid.UUID, date time.Time) (int64, error) {
	// Buscar data de referência do saldo inicial da conta
	var initialBalanceDate *time.Time
	query := `SELECT initial_balance_date FROM accounts WHERE id = $1`
	err := r.db.GetContext(ctx, &initialBalanceDate, query, accountID)
	if err != nil {
		return 0, err
	}

	var balance int64
	var query2 string
	var args []interface{}

	if initialBalanceDate != nil {
		// Se há initial_balance_date, excluir transações de ajuste de saldo inicial na mesma data ou anteriores
		// Isso inclui transações com tag "Saldo Inicial", descrição "Ajuste de saldo inicial" ou tipo adjustment
		query2 = `
			SELECT COALESCE(SUM(
				CASE 
					WHEN type = 'income' THEN amount
					WHEN type = 'adjustment' THEN amount
					WHEN type = 'expense' THEN -amount
					ELSE 0
				END
			), 0)
			FROM transactions
			WHERE account_id = $1 
			  AND date <= $2
			  AND status = 'paid'
			  AND type != 'transfer'
			  AND NOT (
				date <= $3::date
				AND (
					type = 'adjustment'
					OR description = 'Ajuste de saldo inicial'
					OR 'Saldo Inicial' = ANY(tags)
				)
			  )
		`
		args = []interface{}{accountID, date, initialBalanceDate.Format("2006-01-02")}
	} else {
		// Se não há initial_balance_date, incluir todas as transações
		query2 = `
			SELECT COALESCE(SUM(
				CASE 
					WHEN type = 'income' THEN amount
					WHEN type = 'adjustment' THEN amount
					WHEN type = 'expense' THEN -amount
					ELSE 0
				END
			), 0)
			FROM transactions
			WHERE account_id = $1 
			  AND date <= $2
			  AND status = 'paid'
			  AND type != 'transfer'
		`
		args = []interface{}{accountID, date}
	}

	err = r.db.GetContext(ctx, &balance, query2, args...)
	return balance, err
}

// GetAccountInitialBalance retorna o saldo inicial de uma conta (antes de todas as transações)
func (r *transactionRepository) GetAccountInitialBalance(ctx context.Context, accountID uuid.UUID) (int64, error) {
	// Buscar o saldo atual da conta
	var balance int64
	query := `SELECT balance FROM accounts WHERE id = $1`
	err := r.db.GetContext(ctx, &balance, query, accountID)
	if err != nil {
		return 0, err
	}

	// Calcular o saldo de todas as transações
	var transactionBalance int64
	query2 := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN type = 'income' THEN amount
				WHEN type = 'expense' THEN -amount
				ELSE 0
			END
		), 0)
		FROM transactions
		WHERE account_id = $1
	`
	err = r.db.GetContext(ctx, &transactionBalance, query2, accountID)
	if err != nil {
		return 0, err
	}

	// Saldo inicial = saldo atual - saldo das transações
	return balance - transactionBalance, nil
}

// CalculateAccountBalanceFromTransactions calcula o saldo de uma conta baseado no saldo inicial e nas transações
func (r *transactionRepository) CalculateAccountBalanceFromTransactions(ctx context.Context, accountID uuid.UUID) (int64, error) {
	// Buscar saldo inicial e data de referência da conta
	var account struct {
		InitialBalance     *int64     `db:"initial_balance"`
		InitialBalanceDate *time.Time `db:"initial_balance_date"`
	}
	query := `SELECT initial_balance, initial_balance_date FROM accounts WHERE id = $1`
	err := r.db.GetContext(ctx, &account, query, accountID)
	if err != nil {
		return 0, err
	}

	var initialBalance int64
	if account.InitialBalance != nil {
		initialBalance = *account.InitialBalance
	}

	// Calcular o saldo de todas as transações (income + adjustment - expense)
	// Considerar apenas transações pagas e excluir transferências
	// Excluir transações de ajuste que foram criadas na mesma data ou antes do initial_balance_date,
	// pois o initial_balance já representa esse valor
	var transactionBalance int64
	var query2 string
	var args []interface{}

	if account.InitialBalanceDate != nil {
		// Se há initial_balance_date, excluir transações de ajuste de saldo inicial na mesma data ou anteriores
		// Isso inclui transações com tag "Saldo Inicial", descrição "Ajuste de saldo inicial" ou tipo adjustment
		query2 = `
			SELECT COALESCE(SUM(
				CASE 
					WHEN type = 'income' THEN amount
					WHEN type = 'adjustment' THEN amount
					WHEN type = 'expense' THEN -amount
					ELSE 0
				END
			), 0)
			FROM transactions
			WHERE account_id = $1
			  AND status = 'paid'
			  AND type != 'transfer'
			  AND NOT (
				date <= $2::date
				AND (
					type = 'adjustment'
					OR description = 'Ajuste de saldo inicial'
					OR 'Saldo Inicial' = ANY(tags)
				)
			  )
		`
		args = []interface{}{accountID, account.InitialBalanceDate.Format("2006-01-02")}
	} else {
		// Se não há initial_balance_date, incluir todas as transações
		query2 = `
			SELECT COALESCE(SUM(
				CASE 
					WHEN type = 'income' THEN amount
					WHEN type = 'adjustment' THEN amount
					WHEN type = 'expense' THEN -amount
					ELSE 0
				END
			), 0)
			FROM transactions
			WHERE account_id = $1
			  AND status = 'paid'
			  AND type != 'transfer'
		`
		args = []interface{}{accountID}
	}

	err = r.db.GetContext(ctx, &transactionBalance, query2, args...)
	if err != nil {
		return 0, err
	}

	// Saldo total = saldo inicial + transações
	return initialBalance + transactionBalance, nil
}

// FindByStatus busca transações por status
func (r *transactionRepository) FindByStatus(ctx context.Context, userID uuid.UUID, status string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	query := `SELECT id, user_id, account_id, category_id, type, amount, description, date,
	                 status, tags, to_account_id, parent_transaction_id, installment_number,
	                 total_installments, external_id, period_id, reference_month, created_at 
	          FROM transactions 
	          WHERE user_id = $1 AND status = $2
	          ORDER BY date DESC, created_at DESC`
	err := r.db.SelectContext(ctx, &transactions, query, userID, status)
	return transactions, err
}

// FindByTags busca transações que contenham qualquer uma das tags fornecidas
func (r *transactionRepository) FindByTags(ctx context.Context, userID uuid.UUID, tags []string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	query := `SELECT id, user_id, account_id, category_id, type, amount, description, date,
	                 status, tags, to_account_id, parent_transaction_id, installment_number,
	                 total_installments, external_id, period_id, reference_month, created_at 
	          FROM transactions 
	          WHERE user_id = $1 AND tags && $2
	          ORDER BY date DESC, created_at DESC`
	err := r.db.SelectContext(ctx, &transactions, query, userID, tags)
	return transactions, err
}

// FindInstallments busca todas as parcelas de uma transação pai
func (r *transactionRepository) FindInstallments(ctx context.Context, parentTransactionID uuid.UUID) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	query := `SELECT id, user_id, account_id, category_id, type, amount, description, date,
	                 status, tags, to_account_id, parent_transaction_id, installment_number,
	                 total_installments, external_id, period_id, reference_month, created_at 
	          FROM transactions 
	          WHERE parent_transaction_id = $1
	          ORDER BY installment_number ASC`
	err := r.db.SelectContext(ctx, &transactions, query, parentTransactionID)
	return transactions, err
}

// FindByExternalID busca uma transação pelo external_id (para deduplicação)
func (r *transactionRepository) FindByExternalID(ctx context.Context, externalID string) (*model.Transaction, error) {
	var transaction model.Transaction
	query := `SELECT id, user_id, account_id, category_id, type, amount, description, date,
	                 status, tags, to_account_id, parent_transaction_id, installment_number,
	                 total_installments, external_id, period_id, reference_month, created_at 
	          FROM transactions 
	          WHERE external_id = $1`
	err := r.db.GetContext(ctx, &transaction, query, externalID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

// CreateTransfer cria uma transferência entre duas contas usando transação ACID
func (r *transactionRepository) CreateTransfer(ctx context.Context, fromTransaction, toTransaction *model.Transaction, fromAccountID, toAccountID uuid.UUID, amount int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Criar transação de saída
	queryFrom := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date,
		                         status, tags, to_account_id, parent_transaction_id, installment_number,
		                         total_installments, external_id, period_id, reference_month, created_at)
		VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
		        :status, :tags, :to_account_id, :parent_transaction_id, :installment_number,
		        :total_installments, :external_id, :period_id, :reference_month, :created_at)
	`
	_, err = tx.NamedExecContext(ctx, queryFrom, fromTransaction)
	if err != nil {
		return err
	}

	// Criar transação de entrada
	queryTo := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date,
		                         status, tags, to_account_id, parent_transaction_id, installment_number,
		                         total_installments, external_id, period_id, reference_month, created_at)
		VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
		        :status, :tags, :to_account_id, :parent_transaction_id, :installment_number,
		        :total_installments, :external_id, :period_id, :reference_month, :created_at)
	`
	_, err = tx.NamedExecContext(ctx, queryTo, toTransaction)
	if err != nil {
		return err
	}

	// Atualizar saldo da conta origem (subtrair)
	queryUpdateFrom := `UPDATE accounts SET balance = balance - $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, queryUpdateFrom, amount, fromAccountID)
	if err != nil {
		return err
	}

	// Atualizar saldo da conta destino (adicionar)
	queryUpdateTo := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, queryUpdateTo, amount, toAccountID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *transactionRepository) Update(ctx context.Context, transaction *model.Transaction) error {
	query := `
		UPDATE transactions 
		SET account_id = :account_id, category_id = :category_id, type = :type, 
		    amount = :amount, description = :description, date = :date,
		    status = :status, tags = :tags, to_account_id = :to_account_id,
		    parent_transaction_id = :parent_transaction_id, installment_number = :installment_number,
		    total_installments = :total_installments, external_id = :external_id,
		    period_id = :period_id, reference_month = :reference_month
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, transaction)
	return err
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
