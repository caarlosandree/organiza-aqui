package model

import (
	"time"

	"github.com/google/uuid"
)

// RecurrencePattern representa um padrão de recorrência para transações
type RecurrencePattern struct {
	ID                uuid.UUID  `db:"id"`
	UserID            uuid.UUID  `db:"user_id"`
	TransactionID     uuid.UUID  `db:"transaction_id"`
	Frequency         string     `db:"frequency"` // daily, weekly, monthly, yearly
	Interval          int        `db:"interval"`  // a cada X dias/semanas/meses
	EndDate           *time.Time `db:"end_date"`  // NULL = infinito
	LastGeneratedDate *time.Time `db:"last_generated_date"`
	CreatedAt         time.Time  `db:"created_at"`
}
