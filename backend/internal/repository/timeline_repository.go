package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/model"
)

type TimelineRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TimelineFilters) ([]*model.TimelineEvent, error)
	CountByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TimelineFilters) (int, error)
	GetSummary(ctx context.Context, userID uuid.UUID) (*dto.TimelineSummaryDTO, error)
}

type timelineRepository struct {
	db *sqlx.DB
}

func NewTimelineRepository(db *sqlx.DB) TimelineRepository {
	return &timelineRepository{db: db}
}

func (r *timelineRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TimelineFilters) ([]*model.TimelineEvent, error) {
	qb := NewQueryBuilder(`
		SELECT id, user_id, entity_type, entity_id, title, description, event_date, metadata, created_at
		FROM timeline_events
	`).
		WhereEqual("user_id", userID)

	if filters != nil {
		qb.WhereEqual("entity_type", filters.EntityType).
			WhereGreaterOrEqual("event_date", filters.StartDate).
			WhereLessOrEqual("event_date", filters.EndDate).
			OrderBy("event_date DESC, created_at DESC").
			Limit(filters.Limit).
			Offset(filters.Offset)
	} else {
		qb.OrderBy("event_date DESC, created_at DESC")
	}

	query, args := qb.Build()
	var events []*model.TimelineEvent
	err := r.db.SelectContext(ctx, &events, query, args...)
	return events, err
}

func (r *timelineRepository) CountByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TimelineFilters) (int, error) {
	qb := NewQueryBuilder(`SELECT COUNT(*) FROM timeline_events`).
		WhereEqual("user_id", userID)

	if filters != nil {
		qb.WhereEqual("entity_type", filters.EntityType).
			WhereGreaterOrEqual("event_date", filters.StartDate).
			WhereLessOrEqual("event_date", filters.EndDate)
	}

	query, args := qb.Build()
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *timelineRepository) GetSummary(ctx context.Context, userID uuid.UUID) (*dto.TimelineSummaryDTO, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.AddDate(0, 0, 1)

	// Total de eventos
	var totalEvents int
	if err := r.db.GetContext(ctx, &totalEvents, `SELECT COUNT(*) FROM timeline_events WHERE user_id = $1`, userID); err != nil {
		return nil, err
	}

	// Eventos de hoje
	var todayEvents int
	if err := r.db.GetContext(ctx, &todayEvents, `
		SELECT COUNT(*) FROM timeline_events 
		WHERE user_id = $1 AND event_date >= $2 AND event_date < $3
	`, userID, todayStart, todayEnd); err != nil {
		return nil, err
	}

	// Eventos futuros
	var upcomingEvents int
	if err := r.db.GetContext(ctx, &upcomingEvents, `
		SELECT COUNT(*) FROM timeline_events 
		WHERE user_id = $1 AND event_date > $2
	`, userID, now); err != nil {
		return nil, err
	}

	// Por tipo
	typeCounts := make(map[string]int)
	rows, err := r.db.QueryContext(ctx, `
		SELECT entity_type, COUNT(*) 
		FROM timeline_events 
		WHERE user_id = $1 
		GROUP BY entity_type
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entityType string
		var count int
		if err := rows.Scan(&entityType, &count); err != nil {
			return nil, err
		}
		typeCounts[entityType] = count
	}

	return &dto.TimelineSummaryDTO{
		TotalEvents:    totalEvents,
		TodayEvents:    todayEvents,
		UpcomingEvents: upcomingEvents,
		ByType:         typeCounts,
	}, nil
}
