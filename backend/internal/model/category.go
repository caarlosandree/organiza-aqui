package model

import (
	"time"

	"github.com/google/uuid"
)

// Category representa uma categoria financeira
type Category struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	Name      string     `db:"name"`
	ParentID  *uuid.UUID `db:"parent_id"`
	Path      string     `db:"path"` // Materialized Path
	Type      string     `db:"type"` // income, expense
	Color     string     `db:"color"`
	CreatedAt time.Time  `db:"created_at"`
}
