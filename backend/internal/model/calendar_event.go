package model

import (
	"time"

	"github.com/google/uuid"
)

// CalendarEvent representa um evento no calendário
type CalendarEvent struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
	AllDay      bool       `db:"all_day"`
	Location    string     `db:"location"`
	Color       string     `db:"color"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
