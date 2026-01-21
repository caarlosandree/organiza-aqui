package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TimelineEvent representa um evento na timeline unificada
type TimelineEvent struct {
	ID          uuid.UUID       `db:"id"`
	UserID      uuid.UUID       `db:"user_id"`
	EntityType  string          `db:"entity_type"` // transaction, task, calendar_event, note
	EntityID    uuid.UUID       `db:"entity_id"`
	Title       string          `db:"title"`
	Description string          `db:"description"`
	EventDate   time.Time       `db:"event_date"`
	Metadata    json.RawMessage `db:"metadata"` // JSONB
	CreatedAt   time.Time       `db:"created_at"`
}
