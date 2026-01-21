package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
)

type CalendarEventRepository interface {
	Create(ctx context.Context, event *model.CalendarEvent) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.CalendarEvent, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.CalendarEventFilters) ([]*model.CalendarEvent, error)
	Update(ctx context.Context, event *model.CalendarEvent) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type calendarEventRepository struct {
	db *sqlx.DB
}

func NewCalendarEventRepository(db *sqlx.DB) CalendarEventRepository {
	return &calendarEventRepository{db: db}
}

func (r *calendarEventRepository) Create(ctx context.Context, event *model.CalendarEvent) error {
	query := `
		INSERT INTO calendar_events (id, user_id, title, description, start_date, end_date, all_day, location, color, created_at, updated_at)
		VALUES (:id, :user_id, :title, :description, :start_date, :end_date, :all_day, :location, :color, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, event)
	return err
}

func (r *calendarEventRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.CalendarEvent, error) {
	var event model.CalendarEvent
	query := `
		SELECT id, user_id, title, description, start_date, end_date, all_day, location, color, created_at, updated_at
		FROM calendar_events
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *calendarEventRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.CalendarEventFilters) ([]*model.CalendarEvent, error) {
	qb := NewQueryBuilder(`
		SELECT id, user_id, title, description, start_date, end_date, all_day, location, color, created_at, updated_at
		FROM calendar_events
	`).
		WhereEqual("user_id", userID)

	if filters != nil {
		qb.WhereGreaterOrEqual("start_date", filters.StartDate)

		if filters.EndDate != nil {
			qb.WhereCustom("(end_date IS NULL OR end_date <= $1)", *filters.EndDate)
		}

		qb.OrderBy("start_date ASC, created_at ASC").
			Limit(filters.Limit).
			Offset(filters.Offset)
	} else {
		qb.OrderBy("start_date ASC, created_at ASC")
	}

	query, args := qb.Build()
	var events []*model.CalendarEvent
	err := r.db.SelectContext(ctx, &events, query, args...)
	return events, err
}

func (r *calendarEventRepository) Update(ctx context.Context, event *model.CalendarEvent) error {
	query := `
		UPDATE calendar_events
		SET title = :title, description = :description, start_date = :start_date, end_date = :end_date,
		    all_day = :all_day, location = :location, color = :color, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, event)
	return err
}

func (r *calendarEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM calendar_events WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
