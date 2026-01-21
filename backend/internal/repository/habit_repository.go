package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
)

type HabitRepository interface {
	Create(ctx context.Context, habit *model.Habit) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Habit, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Habit, error)
	Update(ctx context.Context, habit *model.Habit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type HabitTrackingRepository interface {
	Create(ctx context.Context, tracking *model.HabitTracking) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.HabitTracking, error)
	FindByHabitID(ctx context.Context, habitID uuid.UUID, startDate, endDate time.Time) ([]*model.HabitTracking, error)
	FindByHabitIDAndDate(ctx context.Context, habitID uuid.UUID, date time.Time) (*model.HabitTracking, error)
	Update(ctx context.Context, tracking *model.HabitTracking) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context, habitID uuid.UUID, startDate, endDate time.Time) (*HabitStats, error)
}

type HabitStats struct {
	TotalDays      int
	CompletedDays  int
	CompletionRate float64
	CurrentStreak  int
	LongestStreak  int
}

type habitRepository struct {
	db *sqlx.DB
}

func NewHabitRepository(db *sqlx.DB) HabitRepository {
	return &habitRepository{db: db}
}

func (r *habitRepository) Create(ctx context.Context, habit *model.Habit) error {
	query := `
		INSERT INTO habits (id, user_id, name, description, color, frequency, target_days, created_at, updated_at)
		VALUES (:id, :user_id, :name, :description, :color, :frequency, :target_days, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, habit)
	return err
}

func (r *habitRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Habit, error) {
	var habit model.Habit
	query := `
		SELECT id, user_id, name, description, color, frequency, target_days, created_at, updated_at
		FROM habits
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &habit, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

func (r *habitRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Habit, error) {
	var habits []*model.Habit
	query := `
		SELECT id, user_id, name, description, color, frequency, target_days, created_at, updated_at
		FROM habits
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &habits, query, userID)
	return habits, err
}

func (r *habitRepository) Update(ctx context.Context, habit *model.Habit) error {
	query := `
		UPDATE habits
		SET name = :name, description = :description, color = :color, frequency = :frequency, 
		    target_days = :target_days, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, habit)
	return err
}

func (r *habitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM habits WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

type habitTrackingRepository struct {
	db *sqlx.DB
}

func NewHabitTrackingRepository(db *sqlx.DB) HabitTrackingRepository {
	return &habitTrackingRepository{db: db}
}

func (r *habitTrackingRepository) Create(ctx context.Context, tracking *model.HabitTracking) error {
	query := `
		INSERT INTO habit_tracking (id, habit_id, date, completed, notes, created_at)
		VALUES (:id, :habit_id, :date, :completed, :notes, :created_at)
		ON CONFLICT (habit_id, date) DO UPDATE SET completed = :completed, notes = :notes
	`
	_, err := r.db.NamedExecContext(ctx, query, tracking)
	return err
}

func (r *habitTrackingRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.HabitTracking, error) {
	var tracking model.HabitTracking
	query := `
		SELECT id, habit_id, date, completed, notes, created_at
		FROM habit_tracking
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &tracking, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tracking, nil
}

func (r *habitTrackingRepository) FindByHabitID(ctx context.Context, habitID uuid.UUID, startDate, endDate time.Time) ([]*model.HabitTracking, error) {
	var trackings []*model.HabitTracking
	query := `
		SELECT id, habit_id, date, completed, notes, created_at
		FROM habit_tracking
		WHERE habit_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`
	err := r.db.SelectContext(ctx, &trackings, query, habitID, startDate, endDate)
	return trackings, err
}

func (r *habitTrackingRepository) FindByHabitIDAndDate(ctx context.Context, habitID uuid.UUID, date time.Time) (*model.HabitTracking, error) {
	var tracking model.HabitTracking
	query := `
		SELECT id, habit_id, date, completed, notes, created_at
		FROM habit_tracking
		WHERE habit_id = $1 AND date = $2
	`
	err := r.db.GetContext(ctx, &tracking, query, habitID, date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tracking, nil
}

func (r *habitTrackingRepository) Update(ctx context.Context, tracking *model.HabitTracking) error {
	query := `
		UPDATE habit_tracking
		SET completed = :completed, notes = :notes
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, tracking)
	return err
}

func (r *habitTrackingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM habit_tracking WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *habitTrackingRepository) GetStats(ctx context.Context, habitID uuid.UUID, startDate, endDate time.Time) (*HabitStats, error) {
	// Total de dias no período
	var totalDays int
	err := r.db.GetContext(ctx, &totalDays, `
		SELECT COUNT(DISTINCT date)
		FROM habit_tracking
		WHERE habit_id = $1 AND date >= $2 AND date <= $3
	`, habitID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Dias completados
	var completedDays int
	err = r.db.GetContext(ctx, &completedDays, `
		SELECT COUNT(*)
		FROM habit_tracking
		WHERE habit_id = $1 AND date >= $2 AND date <= $3 AND completed = true
	`, habitID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Taxa de conclusão
	completionRate := 0.0
	if totalDays > 0 {
		completionRate = float64(completedDays) / float64(totalDays) * 100
	}

	// Calcular streaks (simplificado)
	// Buscar todos os registros ordenados por data
	var trackings []*model.HabitTracking
	err = r.db.SelectContext(ctx, &trackings, `
		SELECT date, completed
		FROM habit_tracking
		WHERE habit_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`, habitID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	currentStreak := 0
	longestStreak := 0
	tempStreak := 0

	for _, t := range trackings {
		if t.Completed {
			tempStreak++
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
		} else {
			if currentStreak == 0 {
				currentStreak = tempStreak
			}
			tempStreak = 0
		}
	}

	if currentStreak == 0 {
		currentStreak = tempStreak
	}

	return &HabitStats{
		TotalDays:      totalDays,
		CompletedDays:  completedDays,
		CompletionRate: completionRate,
		CurrentStreak:  currentStreak,
		LongestStreak:  longestStreak,
	}, nil
}
