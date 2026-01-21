package model

import (
	"time"

	"github.com/google/uuid"
)

// User representa um usuário no sistema
type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
