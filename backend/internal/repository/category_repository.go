package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Category, error)
	FindByUserIDAndType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindChildren(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (id, user_id, name, parent_id, path, type, color, created_at)
		VALUES (:id, :user_id, :name, :parent_id, :path, :type, :color, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	var category model.Category
	query := `SELECT id, user_id, name, parent_id, path, type, color, created_at FROM categories WHERE id = $1`
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Category, error) {
	var categories []*model.Category
	query := `SELECT id, user_id, name, parent_id, path, type, color, created_at FROM categories WHERE user_id = $1 ORDER BY path`
	err := r.db.SelectContext(ctx, &categories, query, userID)
	return categories, err
}

func (r *categoryRepository) FindByUserIDAndType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*model.Category, error) {
	var categories []*model.Category
	query := `SELECT id, user_id, name, parent_id, path, type, color, created_at FROM categories WHERE user_id = $1 AND type = $2 ORDER BY path`
	err := r.db.SelectContext(ctx, &categories, query, userID, categoryType)
	return categories, err
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `
		UPDATE categories 
		SET name = :name, parent_id = :parent_id, path = :path, type = :type, color = :color
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *categoryRepository) FindChildren(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error) {
	var categories []*model.Category
	query := `SELECT id, user_id, name, parent_id, path, type, color, created_at FROM categories WHERE parent_id = $1`
	err := r.db.SelectContext(ctx, &categories, query, parentID)
	return categories, err
}
