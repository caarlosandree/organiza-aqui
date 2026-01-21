package model

import (
	"time"

	"github.com/google/uuid"
)

// Task representa uma tarefa
type Task struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	StatusID    uuid.UUID  `db:"status_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Priority    string     `db:"priority"` // low, medium, high, urgent
	DueDate             *time.Time `db:"due_date"`
	CompletedAt         *time.Time `db:"completed_at"`
	Lexorank            string     `db:"lexorank"` // Para ordenação drag-and-drop
	FinancialAccountID  *uuid.UUID `db:"financial_account_id"`
	FinancialAmount     *int64     `db:"financial_amount"` // em centavos
	FinancialCategoryID *uuid.UUID `db:"financial_category_id"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}
