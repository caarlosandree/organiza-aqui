package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type BankRepository interface {
	Create(ctx context.Context, bank *model.Bank) error
	Update(ctx context.Context, bank *model.Bank) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Bank, error)
	FindByISPB(ctx context.Context, ispb string) (*model.Bank, error)
	FindByCode(ctx context.Context, code int) (*model.Bank, error)
	FindAll(ctx context.Context) ([]*model.Bank, error)
	Upsert(ctx context.Context, bank *model.Bank) error
}

type bankRepository struct {
	db *sqlx.DB
}

func NewBankRepository(db *sqlx.DB) BankRepository {
	return &bankRepository{db: db}
}

func (r *bankRepository) Create(ctx context.Context, bank *model.Bank) error {
	query := `
		INSERT INTO banks (id, ispb, code, name, full_name, created_at, updated_at)
		VALUES (:id, :ispb, :code, :name, :full_name, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, bank)
	return err
}

func (r *bankRepository) Update(ctx context.Context, bank *model.Bank) error {
	query := `
		UPDATE banks 
		SET ispb = :ispb, code = :code, name = :name, full_name = :full_name, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, bank)
	return err
}

func (r *bankRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Bank, error) {
	var bank model.Bank
	query := `SELECT id, ispb, code, name, full_name, created_at, updated_at FROM banks WHERE id = $1`
	err := r.db.GetContext(ctx, &bank, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bank, nil
}

func (r *bankRepository) FindByISPB(ctx context.Context, ispb string) (*model.Bank, error) {
	var bank model.Bank
	query := `SELECT id, ispb, code, name, full_name, created_at, updated_at FROM banks WHERE ispb = $1`
	err := r.db.GetContext(ctx, &bank, query, ispb)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bank, nil
}

func (r *bankRepository) FindByCode(ctx context.Context, code int) (*model.Bank, error) {
	var bank model.Bank
	query := `SELECT id, ispb, code, name, full_name, created_at, updated_at FROM banks WHERE code = $1`
	err := r.db.GetContext(ctx, &bank, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bank, nil
}

func (r *bankRepository) FindAll(ctx context.Context) ([]*model.Bank, error) {
	var banks []*model.Bank
	query := `SELECT id, ispb, code, name, full_name, created_at, updated_at FROM banks ORDER BY name ASC`
	err := r.db.SelectContext(ctx, &banks, query)
	return banks, err
}

func (r *bankRepository) Upsert(ctx context.Context, bank *model.Bank) error {
	existing, err := r.FindByISPB(ctx, bank.ISPB)
	if err != nil {
		return fmt.Errorf("erro ao buscar banco por ISPB: %w", err)
	}

	if existing == nil {
		// Criar novo banco
		if bank.ID == uuid.Nil {
			bank.ID = uuid.New()
		}
		return r.Create(ctx, bank)
	}

	// Atualizar banco existente se houver mudanças
	if existing.Code != bank.Code || existing.Name != bank.Name || existing.FullName != bank.FullName {
		bank.ID = existing.ID
		bank.CreatedAt = existing.CreatedAt
		return r.Update(ctx, bank)
	}

	return nil
}
