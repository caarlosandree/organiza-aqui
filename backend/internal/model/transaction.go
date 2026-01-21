package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Transaction representa uma transação financeira
type Transaction struct {
	ID                 uuid.UUID  `db:"id"`
	UserID             uuid.UUID  `db:"user_id"`
	AccountID          uuid.UUID  `db:"account_id"`
	CategoryID         *uuid.UUID `db:"category_id"`
	Type               string     `db:"type"` // income, expense, transfer, adjustment
	Amount             int64      `db:"amount"` // em centavos
	Description        string     `db:"description"`
	Date               time.Time  `db:"date"`
	Status             string          `db:"status"` // pending, paid, cancelled
	Tags               pq.StringArray  `db:"tags"` // array de tags
	ToAccountID        *uuid.UUID `db:"to_account_id"` // para transferências
	ParentTransactionID *uuid.UUID `db:"parent_transaction_id"` // para parcelas
	InstallmentNumber  *int       `db:"installment_number"` // número da parcela
	TotalInstallments  *int       `db:"total_installments"` // total de parcelas
	ExternalID         *string    `db:"external_id"` // para deduplicação de importações
	PeriodID           *uuid.UUID `db:"period_id"` // ID do período de referência
	ReferenceMonth     *time.Time `db:"reference_month"` // mês de referência (primeiro dia do mês)
	CreatedAt          time.Time  `db:"created_at"`
}
