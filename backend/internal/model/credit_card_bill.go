package model

import (
	"time"

	"github.com/google/uuid"
)

// CreditCardBill representa uma fatura de cartão de crédito
type CreditCardBill struct {
	ID                 uuid.UUID  `db:"id"`
	CreditCardID       uuid.UUID  `db:"credit_card_id"`
	Month              int        `db:"month"`
	Year               int        `db:"year"`
	Status             string     `db:"status"` // open, closed, paid
	ClosingDate        time.Time  `db:"closing_date"`
	DueDate            time.Time  `db:"due_date"`
	PaymentTransactionID *uuid.UUID `db:"payment_transaction_id"` // transação que pagou a fatura
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}
