package model

import (
	"time"

	"github.com/google/uuid"
)

// Account representa uma conta financeira
type Account struct {
	ID                uuid.UUID  `db:"id"`
	UserID            uuid.UUID  `db:"user_id"`
	Name              string     `db:"name"`
	Type              string     `db:"type"` // checking, savings, credit, investment
	Balance           int64      `db:"balance"` // em centavos
	Currency          string     `db:"currency"`
	BankID            uuid.UUID  `db:"bank_id"`
	InitialBalance    *int64     `db:"initial_balance"` // saldo inicial em centavos
	InitialBalanceDate *time.Time `db:"initial_balance_date"` // data de referência do saldo inicial
	CreatedAt         time.Time  `db:"created_at"`
}
