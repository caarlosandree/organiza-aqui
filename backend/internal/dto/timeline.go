package dto

// TimelineEventDTO representa um evento na timeline na camada de apresentação
type TimelineEventDTO struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	EntityType  string                 `json:"entity_type"`
	EntityID    string                 `json:"entity_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	EventDate   string                 `json:"event_date"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   string                 `json:"created_at"`
}

// TimelineFilters representa filtros para listagem de eventos da timeline
type TimelineFilters struct {
	EntityType *string `json:"entity_type,omitempty"`
	StartDate  *string `json:"start_date,omitempty"`
	EndDate    *string `json:"end_date,omitempty"`
	Limit      int     `json:"limit,omitempty"`
	Offset     int     `json:"offset,omitempty"`
}

// TimelineSummaryDTO representa um resumo da timeline
type TimelineSummaryDTO struct {
	TotalEvents    int `json:"total_events"`
	TodayEvents    int `json:"today_events"`
	UpcomingEvents int `json:"upcoming_events"`
	ByType         map[string]int `json:"by_type"`
}
