package dto

// AccountDTO representa uma conta na camada de apresentação
type AccountDTO struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	Name              string  `json:"name"`
	Type              string  `json:"type"`
	Balance           int64   `json:"balance"` // em centavos
	Currency          string  `json:"currency"`
	BankID            string  `json:"bank_id"`
	InitialBalance    *int64  `json:"initial_balance,omitempty"`    // saldo inicial em centavos
	InitialBalanceDate *string `json:"initial_balance_date,omitempty"` // data de referência do saldo inicial
	CreatedAt         string  `json:"created_at"`
}

// CategoryDTO representa uma categoria na camada de apresentação
type CategoryDTO struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	ParentID  *string `json:"parent_id,omitempty"`
	Path      string  `json:"path"`
	Type      string  `json:"type"`
	Color     string  `json:"color"`
	CreatedAt string  `json:"created_at"`
	Children  []CategoryDTO `json:"children,omitempty"` // para hierarquia
}

// TransactionDTO representa uma transação na camada de apresentação
type TransactionDTO struct {
	ID                 string   `json:"id"`
	UserID             string   `json:"user_id"`
	AccountID          string   `json:"account_id"`
	CategoryID         *string  `json:"category_id,omitempty"`
	Type               string   `json:"type"`
	Amount             int64    `json:"amount"` // em centavos
	Description        string   `json:"description"`
	Date               string   `json:"date"`
	ReferenceMonth     *string  `json:"reference_month,omitempty"` // formato YYYY-MM
	Status             string   `json:"status"` // pending, paid, cancelled
	Tags               []string `json:"tags,omitempty"`
	ToAccountID        *string  `json:"to_account_id,omitempty"` // para transferências
	ParentTransactionID *string `json:"parent_transaction_id,omitempty"` // para parcelas
	InstallmentNumber  *int     `json:"installment_number,omitempty"`
	TotalInstallments  *int     `json:"total_installments,omitempty"`
	PeriodID           *string  `json:"period_id,omitempty"`
	CreatedAt          string   `json:"created_at"`
}

// CreateAccountRequest representa a requisição de criação de conta
type CreateAccountRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Type     string `json:"type" validate:"required,oneof=checking savings credit investment"`
	Currency string `json:"currency" validate:"required,len=3"`
	BankID   string `json:"bank_id" validate:"required,uuid"`
}

// UpdateAccountRequest representa a requisição de atualização de conta
type UpdateAccountRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Type     string `json:"type" validate:"required,oneof=checking savings credit investment"`
	Currency string `json:"currency" validate:"required,len=3"`
	BankID   string `json:"bank_id" validate:"required,uuid"`
}

// UpdateInitialBalanceRequest representa a requisição de atualização de saldo inicial
type UpdateInitialBalanceRequest struct {
	Balance int64  `json:"balance" validate:"required"` // saldo em centavos
	Date    string `json:"date" validate:"required"`   // data de referência no formato "2006-01-02"
}

// CreateCategoryRequest representa a requisição de criação de categoria
type CreateCategoryRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=255"`
	ParentID *string `json:"parent_id,omitempty"`
	Type     string  `json:"type" validate:"required,oneof=income expense"`
	Color    string  `json:"color" validate:"required"`
}

// UpdateCategoryRequest representa a requisição de atualização de categoria
type UpdateCategoryRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=255"`
	ParentID *string `json:"parent_id,omitempty"`
	Type     string  `json:"type" validate:"required,oneof=income expense"`
	Color    string  `json:"color" validate:"required"`
}

// CreateTransactionRequest representa a requisição de criação de transação
type CreateTransactionRequest struct {
	AccountID          string   `json:"account_id" validate:"required,uuid"`
	CategoryID        *string  `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Type              string   `json:"type" validate:"required,oneof=income expense transfer adjustment"`
	Amount            int64    `json:"amount" validate:"required,gt=0"`
	Description       string   `json:"description" validate:"max=1000"`
	Date              string   `json:"date" validate:"required"`
	ReferenceMonth    *string  `json:"reference_month,omitempty" validate:"omitempty"` // formato YYYY-MM
	Status            string   `json:"status,omitempty" validate:"omitempty,oneof=pending paid cancelled"`
	Tags              []string `json:"tags,omitempty"`
	ToAccountID       *string  `json:"to_account_id,omitempty" validate:"omitempty,uuid"` // para transferências
	TotalInstallments *int     `json:"total_installments,omitempty" validate:"omitempty,min=1,max=999"` // para parcelamento
}

// UpdateTransactionRequest representa a requisição de atualização de transação
type UpdateTransactionRequest struct {
	AccountID      string   `json:"account_id" validate:"required,uuid"`
	CategoryID     *string  `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Type           string   `json:"type" validate:"required,oneof=income expense transfer adjustment"`
	Amount         int64    `json:"amount" validate:"required,gt=0"`
	Description    string   `json:"description" validate:"max=1000"`
	Date           string   `json:"date" validate:"required"`
	ReferenceMonth *string  `json:"reference_month,omitempty" validate:"omitempty"` // formato YYYY-MM
	Status         string   `json:"status,omitempty" validate:"omitempty,oneof=pending paid cancelled"`
	Tags           []string `json:"tags,omitempty"`
	ToAccountID    *string  `json:"to_account_id,omitempty" validate:"omitempty,uuid"`
}

// TransactionFilters representa filtros para listagem de transações
type TransactionFilters struct {
	AccountID          *string  `json:"account_id,omitempty"`
	CategoryID         *string  `json:"category_id,omitempty"`
	Type               *string  `json:"type,omitempty"`
	Status             *string  `json:"status,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	StartDate          *string  `json:"start_date,omitempty"`
	EndDate            *string  `json:"end_date,omitempty"`
	MinAmount          *int64   `json:"min_amount,omitempty"` // em centavos
	MaxAmount          *int64   `json:"max_amount,omitempty"` // em centavos
	ParentTransactionID *string  `json:"parent_transaction_id,omitempty"` // para buscar parcelas
	Limit              int      `json:"limit,omitempty"`
	Offset             int      `json:"offset,omitempty"`
}

// StatementResponse representa a resposta de um extrato
type StatementResponse struct {
	AccountID      string                `json:"account_id"`
	AccountName    string                `json:"account_name"`
	StartDate      string                `json:"start_date"`
	EndDate        string                `json:"end_date"`
	InitialBalance int64                 `json:"initial_balance"` // em centavos
	FinalBalance   int64                 `json:"final_balance"`   // em centavos
	TotalIncome    int64                 `json:"total_income"`    // em centavos
	TotalExpense   int64                 `json:"total_expense"`   // em centavos
	Transactions   []*TransactionDTO     `json:"transactions"`
	Summary        StatementSummaryDTO   `json:"summary"`
}

// StatementSummaryDTO representa o resumo do extrato
type StatementSummaryDTO struct {
	TransactionCount int   `json:"transaction_count"`
	IncomeCount      int   `json:"income_count"`
	ExpenseCount     int   `json:"expense_count"`
	AverageIncome    int64 `json:"average_income"`    // em centavos
	AverageExpense   int64 `json:"average_expense"`   // em centavos
}

// IncomeExpenseByPeriodResponse representa receitas e despesas por período
type IncomeExpenseByPeriodResponse struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	TotalIncome int64  `json:"total_income"`  // em centavos
	TotalExpense int64 `json:"total_expense"`  // em centavos
	Balance     int64  `json:"balance"`        // em centavos
}

// CategoryBreakdownDTO representa o breakdown por categoria
type CategoryBreakdownDTO struct {
	CategoryID       string `json:"category_id"`
	CategoryName     string `json:"category_name,omitempty"`
	TotalAmount      int64  `json:"total_amount"`      // em centavos
	TransactionCount int    `json:"transaction_count"`
	Percentage       float64 `json:"percentage,omitempty"`
}

// MonthlyTrendDTO representa a tendência mensal
type MonthlyTrendDTO struct {
	Month   string `json:"month"`   // formato: "2006-01"
	Income  int64  `json:"income"`  // em centavos
	Expense int64  `json:"expense"` // em centavos
}

// CreditCardDTO representa um cartão de crédito na camada de apresentação
type CreditCardDTO struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	AccountID   string `json:"account_id"`
	LimitAmount int64  `json:"limit_amount"` // em centavos
	ClosingDay  int    `json:"closing_day"`
	DueDay      int    `json:"due_day"`
	Color       string `json:"color"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateCreditCardRequest representa a requisição de criação de cartão de crédito
type CreateCreditCardRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	AccountID   string `json:"account_id" validate:"required,uuid"`
	LimitAmount int64  `json:"limit_amount" validate:"required,gt=0"`
	ClosingDay  int    `json:"closing_day" validate:"required,min=1,max=31"`
	DueDay      int    `json:"due_day" validate:"required,min=1,max=31"`
	Color       string `json:"color" validate:"required"`
}

// UpdateCreditCardRequest representa a requisição de atualização de cartão de crédito
type UpdateCreditCardRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	AccountID   string `json:"account_id" validate:"required,uuid"`
	LimitAmount int64  `json:"limit_amount" validate:"required,gt=0"`
	ClosingDay  int    `json:"closing_day" validate:"required,min=1,max=31"`
	DueDay      int    `json:"due_day" validate:"required,min=1,max=31"`
	Color       string `json:"color" validate:"required"`
}

// CreditCardBillDTO representa uma fatura de cartão de crédito
type CreditCardBillDTO struct {
	ID                 string  `json:"id"`
	CreditCardID       string  `json:"credit_card_id"`
	Month              int     `json:"month"`
	Year               int     `json:"year"`
	Status             string  `json:"status"` // open, closed, paid
	ClosingDate        string  `json:"closing_date"`
	DueDate            string  `json:"due_date"`
	TotalAmount        int64   `json:"total_amount"` // calculado dinamicamente
	PaymentTransactionID *string `json:"payment_transaction_id,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

// PayBillRequest representa a requisição de pagamento de fatura
type PayBillRequest struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
	Date      string `json:"date" validate:"required"`
}

// UpdateInstallmentRequest representa a requisição de atualização de parcelas
type UpdateInstallmentRequest struct {
	Amount      int64    `json:"amount" validate:"required,gt=0"`
	Description string   `json:"description" validate:"max=1000"`
	Date        string   `json:"date" validate:"required"`
	Status      string   `json:"status,omitempty" validate:"omitempty,oneof=pending paid cancelled"`
	Tags        []string `json:"tags,omitempty"`
	Scope       string   `json:"scope" validate:"required,oneof=this this_and_future all"` // escopo da edição
}

// ImportOFXRequest representa a requisição de importação OFX
type ImportOFXRequest struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
	File      []byte `json:"file" validate:"required"`
}

// ImportCSVRequest representa a requisição de importação CSV
type ImportCSVRequest struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
	File      []byte `json:"file" validate:"required"`
	Delimiter string `json:"delimiter,omitempty"` // padrão: ","
}

// ImportPreviewResponse representa o preview da importação
type ImportPreviewResponse struct {
	TotalTransactions int      `json:"total_transactions"`
	Duplicates        int      `json:"duplicates"`
	NewTransactions   int      `json:"new_transactions"`
	Transactions      []*TransactionDTO `json:"transactions"`
}

// ImportResponse representa a resposta da importação
type ImportResponse struct {
	TotalProcessed int `json:"total_processed"`
	Duplicates     int `json:"duplicates"`
	Created        int `json:"created"`
	Errors         int `json:"errors"`
}

// PatrimonioLiquidoResponse representa o patrimônio líquido
type PatrimonioLiquidoResponse struct {
	TotalPatrimonio int64 `json:"total_patrimonio"` // em centavos
	TotalContas     int64 `json:"total_contas"`     // em centavos
	TotalFaturas    int64 `json:"total_faturas"`    // em centavos (faturas abertas)
	Contas          []AccountDTO `json:"contas,omitempty"`
}

// CalendarioVencimentosResponse representa o calendário de vencimentos
type CalendarioVencimentosResponse struct {
	Items []VencimentoItemDTO `json:"items"`
}

// VencimentoItemDTO representa um item de vencimento
type VencimentoItemDTO struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // "transaction" ou "bill"
	Description string `json:"description"`
	Amount      int64  `json:"amount"` // em centavos
	Date        string `json:"date"`
	Status      string `json:"status"`
}

// GastosPorTagResponse representa gastos agrupados por tag
type GastosPorTagResponse struct {
	Tags []TagBreakdownDTO `json:"tags"`
}

// TagBreakdownDTO representa breakdown por tag
type TagBreakdownDTO struct {
	Tag              string  `json:"tag"`
	TotalAmount      int64   `json:"total_amount"`      // em centavos
	TransactionCount int     `json:"transaction_count"`
	Percentage       float64 `json:"percentage,omitempty"`
}

// ProjecaoFaturaResponse representa a projeção de fatura
type ProjecaoFaturaResponse struct {
	CreditCardID string `json:"credit_card_id"`
	Month        int    `json:"month"`
	Year         int    `json:"year"`
	ProjectedAmount int64 `json:"projected_amount"` // em centavos
	ClosingDate  string `json:"closing_date"`
	DueDate      string `json:"due_date"`
}

// TransactionPeriodDTO representa um período de transações na camada de apresentação
type TransactionPeriodDTO struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	AccountID  string `json:"account_id"`
	AccountName *string `json:"account_name,omitempty"`
	PeriodType string `json:"period_type"` // "bank" ou "credit_card"
	Year       int    `json:"year"`
	Month      int    `json:"month"`
	Status     string `json:"status"` // "open", "closed", "archived"
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// PeriodWithTransactionsDTO representa um período com suas transações e estatísticas
type PeriodWithTransactionsDTO struct {
	Period              TransactionPeriodDTO `json:"period"`
	Transactions        []*TransactionDTO    `json:"transactions"`
	TotalIncome         int64                `json:"total_income"`         // receitas bancárias
	TotalBankExpense    int64                `json:"total_bank_expense"`   // despesas bancárias
	TotalCreditCardExpense int64             `json:"total_credit_card_expense"` // faturas cartão
	Balance             int64                `json:"balance"`             // receitas - despesas bancárias (sem cartão)
}

// TransactionPeriodFilters representa filtros para listagem de períodos
type TransactionPeriodFilters struct {
	AccountID    *string `json:"account_id,omitempty"`
	PeriodType   *string `json:"period_type,omitempty"` // "bank" ou "credit_card"
	Year         *int    `json:"year,omitempty"`
	Month        *int    `json:"month,omitempty"`
	Status       *string `json:"status,omitempty"` // "open", "closed", "archived"
	ReferenceMonth *string `json:"reference_month,omitempty"` // formato YYYY-MM
}
