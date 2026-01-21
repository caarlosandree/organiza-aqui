package model

import (
	"time"

	"github.com/google/uuid"
)

// TransactionPeriod representa um período de transações (mês/ano)
type TransactionPeriod struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	AccountID  uuid.UUID `db:"account_id"`
	PeriodType string    `db:"period_type"` // "bank" ou "credit_card"
	Year       int       `db:"year"`
	Month      int       `db:"month"`
	Status     string    `db:"status"` // "open", "closed", "archived"
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
