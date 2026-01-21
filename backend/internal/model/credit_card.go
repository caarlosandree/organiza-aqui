package model

import (
	"time"

	"github.com/google/uuid"
)

// CreditCard representa um cartão de crédito
type CreditCard struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	AccountID   uuid.UUID `db:"account_id"`
	LimitAmount int64     `db:"limit_amount"` // em centavos
	ClosingDay  int       `db:"closing_day"`  // dia do fechamento (1-31)
	DueDay      int       `db:"due_day"`      // dia do vencimento (1-31)
	Color       string    `db:"color"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
