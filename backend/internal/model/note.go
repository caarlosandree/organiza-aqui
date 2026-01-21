package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Note representa uma anotação
type Note struct {
	ID        uuid.UUID      `db:"id"`
	UserID    uuid.UUID      `db:"user_id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	Tags      pq.StringArray `db:"tags"` // Array de strings do PostgreSQL
	IsPinned  bool           `db:"is_pinned"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
