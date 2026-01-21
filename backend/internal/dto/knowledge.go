package dto

// CalendarEventDTO representa um evento do calendário na camada de apresentação
type CalendarEventDTO struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	AllDay      bool    `json:"all_day"`
	Location    string  `json:"location,omitempty"`
	Color       string  `json:"color"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateCalendarEventRequest representa a requisição de criação de evento
type CreateCalendarEventRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=5000"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date,omitempty"`
	AllDay      bool    `json:"all_day"`
	Location    string  `json:"location" validate:"max=255"`
	Color       string  `json:"color" validate:"required"`
}

// UpdateCalendarEventRequest representa a requisição de atualização de evento
type UpdateCalendarEventRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=5000"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date,omitempty"`
	AllDay      bool    `json:"all_day"`
	Location    string  `json:"location" validate:"max=255"`
	Color       string  `json:"color" validate:"required"`
}

// CalendarEventFilters representa filtros para listagem de eventos
type CalendarEventFilters struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Limit     int     `json:"limit,omitempty"`
	Offset    int     `json:"offset,omitempty"`
}

// NoteDTO representa uma anotação na camada de apresentação
type NoteDTO struct {
	ID        string   `json:"id"`
	UserID    string   `json:"user_id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	IsPinned  bool     `json:"is_pinned"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// CreateNoteRequest representa a requisição de criação de anotação
type CreateNoteRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=255"`
	Content  string   `json:"content" validate:"required"`
	Tags     []string `json:"tags"`
	IsPinned bool     `json:"is_pinned"`
}

// UpdateNoteRequest representa a requisição de atualização de anotação
type UpdateNoteRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=255"`
	Content  string   `json:"content" validate:"required"`
	Tags     []string `json:"tags"`
	IsPinned bool     `json:"is_pinned"`
}

// NoteFilters representa filtros para listagem de anotações
type NoteFilters struct {
	Tag      *string `json:"tag,omitempty"`
	IsPinned *bool   `json:"is_pinned,omitempty"`
	Limit    int     `json:"limit,omitempty"`
	Offset   int     `json:"offset,omitempty"`
}

// HabitDTO representa um hábito na camada de apresentação
type HabitDTO struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Frequency   string `json:"frequency"`
	TargetDays  int    `json:"target_days"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateHabitRequest representa a requisição de criação de hábito
type CreateHabitRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=5000"`
	Color       string `json:"color" validate:"required"`
	Frequency   string `json:"frequency" validate:"required,oneof=daily weekly monthly"`
	TargetDays  int    `json:"target_days" validate:"min=1"`
}

// UpdateHabitRequest representa a requisição de atualização de hábito
type UpdateHabitRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=5000"`
	Color       string `json:"color" validate:"required"`
	Frequency   string `json:"frequency" validate:"required,oneof=daily weekly monthly"`
	TargetDays  int    `json:"target_days" validate:"min=1"`
}

// HabitTrackingDTO representa um registro de tracking de hábito
type HabitTrackingDTO struct {
	ID        string `json:"id"`
	HabitID   string `json:"habit_id"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
	Notes     string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
}

// CreateHabitTrackingRequest representa a requisição de criação de tracking
type CreateHabitTrackingRequest struct {
	HabitID   string `json:"habit_id" validate:"required,uuid"`
	Date      string `json:"date" validate:"required"`
	Completed bool   `json:"completed"`
	Notes     string `json:"notes" validate:"max=1000"`
}

// UpdateHabitTrackingRequest representa a requisição de atualização de tracking
type UpdateHabitTrackingRequest struct {
	Completed bool   `json:"completed"`
	Notes     string `json:"notes" validate:"max=1000"`
}

// HabitStatsDTO representa estatísticas de um hábito
type HabitStatsDTO struct {
	HabitID         string  `json:"habit_id"`
	TotalDays        int     `json:"total_days"`
	CompletedDays    int     `json:"completed_days"`
	CompletionRate   float64 `json:"completion_rate"` // Percentual
	CurrentStreak    int     `json:"current_streak"`  // Sequência atual
	LongestStreak    int     `json:"longest_streak"`  // Maior sequência
}
