package model

import (
	"time"

	"github.com/google/uuid"
)

// Habit representa um hábito
type Habit struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Color       string    `db:"color"`
	Frequency   string    `db:"frequency"` // daily, weekly, monthly
	TargetDays  int       `db:"target_days"` // Quantos dias por período
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// HabitTracking representa um registro de execução de hábito
type HabitTracking struct {
	ID        uuid.UUID  `db:"id"`
	HabitID   uuid.UUID  `db:"habit_id"`
	Date      time.Time  `db:"date"`
	Completed bool       `db:"completed"`
	Notes     string     `db:"notes"`
	CreatedAt time.Time  `db:"created_at"`
}
