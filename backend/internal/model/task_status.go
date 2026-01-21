package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus representa um status de tarefa (ex: "To Do", "In Progress", "Done")
type TaskStatus struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	Name       string    `db:"name"`
	Color      string    `db:"color"`
	OrderIndex int       `db:"order_index"`
	IsDefault  bool      `db:"is_default"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
