package dto

// RecurrencePatternDTO representa um padrão de recorrência na camada de apresentação
type RecurrencePatternDTO struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	TransactionID     string  `json:"transaction_id"`
	Frequency         string  `json:"frequency"`
	Interval          int     `json:"interval"`
	EndDate           *string `json:"end_date,omitempty"`
	LastGeneratedDate *string `json:"last_generated_date,omitempty"`
	CreatedAt         string  `json:"created_at"`
}

// CreateRecurrencePatternRequest representa a requisição de criação de padrão de recorrência
type CreateRecurrencePatternRequest struct {
	TransactionID string  `json:"transaction_id" validate:"required,uuid"`
	Frequency     string  `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	Interval      int     `json:"interval" validate:"required,min=1"`
	EndDate       *string `json:"end_date,omitempty" validate:"omitempty"`
}

// UpdateRecurrencePatternRequest representa a requisição de atualização de padrão de recorrência
type UpdateRecurrencePatternRequest struct {
	Frequency string  `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	Interval  int     `json:"interval" validate:"required,min=1"`
	EndDate   *string `json:"end_date,omitempty" validate:"omitempty"`
}

// GenerateRecurrenceRequest representa a requisição para gerar transações recorrentes
type GenerateRecurrenceRequest struct {
	UntilDate string `json:"until_date" validate:"required"` // Gerar transações até esta data
}
