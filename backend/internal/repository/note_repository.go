package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Note, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.NoteFilters) ([]*model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type noteRepository struct {
	db *sqlx.DB
}

func NewNoteRepository(db *sqlx.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(ctx context.Context, note *model.Note) error {
	query := `
		INSERT INTO notes (id, user_id, title, content, tags, is_pinned, created_at, updated_at)
		VALUES (:id, :user_id, :title, :content, :tags, :is_pinned, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, note)
	return err
}

func (r *noteRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	var note model.Note
	query := `
		SELECT id, user_id, title, content, tags, is_pinned, created_at, updated_at
		FROM notes
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &note, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.NoteFilters) ([]*model.Note, error) {
	qb := NewQueryBuilder(`
		SELECT id, user_id, title, content, tags, is_pinned, created_at, updated_at
		FROM notes
	`).
		WhereEqual("user_id", userID)

	if filters != nil {
		if filters.Tag != nil {
			qb.WhereArrayContains("tags", *filters.Tag)
		}
		qb.WhereEqual("is_pinned", filters.IsPinned).
			OrderBy("is_pinned DESC, created_at DESC").
			Limit(filters.Limit).
			Offset(filters.Offset)
	} else {
		qb.OrderBy("is_pinned DESC, created_at DESC")
	}

	query, args := qb.Build()
	var notes []*model.Note
	err := r.db.SelectContext(ctx, &notes, query, args...)
	return notes, err
}

func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	query := `
		UPDATE notes
		SET title = :title, content = :content, tags = :tags, is_pinned = :is_pinned, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, note)
	return err
}

func (r *noteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
