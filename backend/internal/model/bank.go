package model

import (
	"time"

	"github.com/google/uuid"
)

// Bank representa um banco
type Bank struct {
	ID        uuid.UUID `db:"id"`
	ISPB      string    `db:"ispb"`
	Code      int       `db:"code"`
	Name      string    `db:"name"`
	FullName  string    `db:"full_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
