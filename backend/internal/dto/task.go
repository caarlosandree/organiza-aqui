package dto

// TaskStatusDTO representa um status de tarefa na camada de apresentação
type TaskStatusDTO struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"order_index"`
	IsDefault  bool   `json:"is_default"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// TaskDTO representa uma tarefa na camada de apresentação
type TaskDTO struct {
	ID                  string  `json:"id"`
	UserID              string  `json:"user_id"`
	StatusID            string  `json:"status_id"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Priority            string  `json:"priority"`
	DueDate             *string `json:"due_date,omitempty"`
	CompletedAt         *string `json:"completed_at,omitempty"`
	Lexorank            string  `json:"lexorank"`
	FinancialAccountID  *string `json:"financial_account_id,omitempty"`
	FinancialAmount     *int64  `json:"financial_amount,omitempty"` // em centavos
	FinancialCategoryID *string `json:"financial_category_id,omitempty"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

// CreateTaskStatusRequest representa a requisição de criação de status
type CreateTaskStatusRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=255"`
	Color      string `json:"color" validate:"required"`
	OrderIndex int    `json:"order_index" validate:"min=0"`
	IsDefault  bool   `json:"is_default"`
}

// UpdateTaskStatusRequest representa a requisição de atualização de status
type UpdateTaskStatusRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=255"`
	Color      string `json:"color" validate:"required"`
	OrderIndex int    `json:"order_index" validate:"min=0"`
	IsDefault  bool   `json:"is_default"`
}

// CreateTaskRequest representa a requisição de criação de tarefa
type CreateTaskRequest struct {
	StatusID            string  `json:"status_id" validate:"required,uuid"`
	Title               string  `json:"title" validate:"required,min=1,max=255"`
	Description         string  `json:"description" validate:"max=5000"`
	Priority            string  `json:"priority" validate:"required,oneof=low medium high urgent"`
	DueDate             *string `json:"due_date,omitempty"`
	FinancialAccountID  *string `json:"financial_account_id,omitempty" validate:"omitempty,uuid"`
	FinancialAmount     *int64  `json:"financial_amount,omitempty" validate:"omitempty,gt=0"`
	FinancialCategoryID *string `json:"financial_category_id,omitempty" validate:"omitempty,uuid"`
}

// UpdateTaskRequest representa a requisição de atualização de tarefa
type UpdateTaskRequest struct {
	StatusID            string  `json:"status_id" validate:"required,uuid"`
	Title               string  `json:"title" validate:"required,min=1,max=255"`
	Description         string  `json:"description" validate:"max=5000"`
	Priority            string  `json:"priority" validate:"required,oneof=low medium high urgent"`
	DueDate             *string `json:"due_date,omitempty"`
	FinancialAccountID  *string `json:"financial_account_id,omitempty" validate:"omitempty,uuid"`
	FinancialAmount     *int64  `json:"financial_amount,omitempty" validate:"omitempty,gt=0"`
	FinancialCategoryID *string `json:"financial_category_id,omitempty" validate:"omitempty,uuid"`
}

// ReorderTasksRequest representa a requisição de reordenação de tarefas
type ReorderTasksRequest struct {
	TaskID   string `json:"task_id" validate:"required,uuid"`
	StatusID string `json:"status_id" validate:"required,uuid"`
	AfterID  *string `json:"after_id,omitempty" validate:"omitempty,uuid"` // ID da tarefa após a qual inserir (null = início)
}

// TaskFilters representa filtros para listagem de tarefas
type TaskFilters struct {
	StatusID  *string `json:"status_id,omitempty"`
	Priority  *string `json:"priority,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
	Limit     int     `json:"limit,omitempty"`
	Offset    int     `json:"offset,omitempty"`
}
