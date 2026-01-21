package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/internal/repository"
)

type AnalyticsService interface {
	GetIncomeExpenseByPeriod(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*dto.IncomeExpenseByPeriodResponse, error)
	GetCategoryBreakdown(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, transactionType string) ([]*dto.CategoryBreakdownDTO, error)
	GetMonthlyTrend(ctx context.Context, userID uuid.UUID, months int) ([]*dto.MonthlyTrendDTO, error)
	GetPatrimonioLiquido(ctx context.Context, userID uuid.UUID) (*dto.PatrimonioLiquidoResponse, error)
	GetCalendarioVencimentos(ctx context.Context, userID uuid.UUID, days int) (*dto.CalendarioVencimentosResponse, error)
	GetGastosPorTag(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*dto.GastosPorTagResponse, error)
}

type analyticsService struct {
	transactionRepo  repository.TransactionRepository
	accountRepo      repository.AccountRepository
	creditCardRepo   repository.CreditCardRepository
	creditCardBillRepo repository.CreditCardBillRepository
}

func NewAnalyticsService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	creditCardRepo repository.CreditCardRepository,
	creditCardBillRepo repository.CreditCardBillRepository,
) AnalyticsService {
	return &analyticsService{
		transactionRepo:  transactionRepo,
		accountRepo:      accountRepo,
		creditCardRepo:   creditCardRepo,
		creditCardBillRepo: creditCardBillRepo,
	}
}

func (s *analyticsService) GetIncomeExpenseByPeriod(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*dto.IncomeExpenseByPeriodResponse, error) {
	filters := &dto.TransactionFilters{
		StartDate: func() *string { s := startDate.Format("2006-01-02"); return &s }(),
		EndDate:   func() *string { s := endDate.Format("2006-01-02"); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	var totalIncome, totalExpense int64
	for _, t := range transactions {
		if t.Type == "income" {
			totalIncome += t.Amount
		} else if t.Type == "expense" {
			totalExpense += t.Amount
		}
	}

	return &dto.IncomeExpenseByPeriodResponse{
		StartDate:    startDate.Format("2006-01-02"),
		EndDate:      endDate.Format("2006-01-02"),
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
	}, nil
}

func (s *analyticsService) GetCategoryBreakdown(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, transactionType string) ([]*dto.CategoryBreakdownDTO, error) {
	typeFilter := transactionType
	filters := &dto.TransactionFilters{
		Type:      &typeFilter,
		StartDate: func() *string { s := startDate.Format("2006-01-02"); return &s }(),
		EndDate:   func() *string { s := endDate.Format("2006-01-02"); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Agrupar por categoria
	categoryMap := make(map[uuid.UUID]*dto.CategoryBreakdownDTO)
	for _, t := range transactions {
		if t.CategoryID == nil {
			continue
		}

		if breakdown, exists := categoryMap[*t.CategoryID]; exists {
			breakdown.TotalAmount += t.Amount
			breakdown.TransactionCount++
		} else {
			categoryMap[*t.CategoryID] = &dto.CategoryBreakdownDTO{
				CategoryID:       t.CategoryID.String(),
				TotalAmount:      t.Amount,
				TransactionCount: 1,
			}
		}
	}

	// Converter mapa para slice
	breakdowns := make([]*dto.CategoryBreakdownDTO, 0, len(categoryMap))
	for _, breakdown := range categoryMap {
		breakdowns = append(breakdowns, breakdown)
	}

	return breakdowns, nil
}

func (s *analyticsService) GetMonthlyTrend(ctx context.Context, userID uuid.UUID, months int) ([]*dto.MonthlyTrendDTO, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startDate = startDate.AddDate(0, -months+1, 0)

	endDate := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())

	filters := &dto.TransactionFilters{
		StartDate: func() *string { s := startDate.Format("2006-01-02"); return &s }(),
		EndDate:   func() *string { s := endDate.Format("2006-01-02"); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Agrupar por mês
	monthMap := make(map[string]*dto.MonthlyTrendDTO)
	for _, t := range transactions {
		monthKey := t.Date.Format("2006-01")
		if trend, exists := monthMap[monthKey]; exists {
			if t.Type == "income" {
				trend.Income += t.Amount
			} else if t.Type == "expense" {
				trend.Expense += t.Amount
			}
		} else {
			trend := &dto.MonthlyTrendDTO{
				Month:   monthKey,
				Income:  0,
				Expense: 0,
			}
			if t.Type == "income" {
				trend.Income = t.Amount
			} else if t.Type == "expense" {
				trend.Expense = t.Amount
			}
			monthMap[monthKey] = trend
		}
	}

	// Converter para slice ordenado
	trends := make([]*dto.MonthlyTrendDTO, 0, len(monthMap))
	for month := startDate; !month.After(endDate); month = month.AddDate(0, 1, 0) {
		monthKey := month.Format("2006-01")
		if trend, exists := monthMap[monthKey]; exists {
			trends = append(trends, trend)
		} else {
			trends = append(trends, &dto.MonthlyTrendDTO{
				Month:   monthKey,
				Income:  0,
				Expense: 0,
			})
		}
	}

	return trends, nil
}

// GetPatrimonioLiquido calcula o patrimônio líquido (soma de saldos de todas contas - faturas abertas de cartões)
func (s *analyticsService) GetPatrimonioLiquido(ctx context.Context, userID uuid.UUID) (*dto.PatrimonioLiquidoResponse, error) {
	// Buscar todas as contas do usuário
	accounts, err := s.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalContas int64
	accountDTOs := make([]dto.AccountDTO, len(accounts))
	for i, account := range accounts {
		totalContas += account.Balance
		accountDTOs[i] = dto.AccountDTO{
			ID:        account.ID.String(),
			UserID:    account.UserID.String(),
			Name:      account.Name,
			Type:      account.Type,
			Balance:   account.Balance,
			Currency:  account.Currency,
			CreatedAt: account.CreatedAt.Format(time.RFC3339),
		}
	}

	// Buscar todas as faturas abertas de cartões
	creditCards, err := s.creditCardRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalFaturas int64
	for _, cc := range creditCards {
		// Buscar fatura aberta
		openBill, err := s.creditCardRepo.FindOpenBill(ctx, cc.ID)
		if err != nil {
			continue
		}
		if openBill != nil {
			// Calcular total da fatura
			total, err := s.creditCardBillRepo.CalculateBillTotal(ctx, openBill.ID)
			if err == nil {
				totalFaturas += total
			}
		}
	}

	totalPatrimonio := totalContas - totalFaturas

	return &dto.PatrimonioLiquidoResponse{
		TotalPatrimonio: totalPatrimonio,
		TotalContas:     totalContas,
		TotalFaturas:    totalFaturas,
		Contas:          accountDTOs,
	}, nil
}

// GetCalendarioVencimentos retorna transações e faturas com status 'pending' ordenadas por data
func (s *analyticsService) GetCalendarioVencimentos(ctx context.Context, userID uuid.UUID, days int) (*dto.CalendarioVencimentosResponse, error) {
	// Data limite
	endDate := time.Now().AddDate(0, 0, days)

	// Buscar transações pendentes
	pendingStatus := "pending"
	filters := &dto.TransactionFilters{
		Status:   &pendingStatus,
		EndDate:  func() *string { s := endDate.Format("2006-01-02"); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	var items []dto.VencimentoItemDTO

	// Adicionar transações pendentes
	for _, t := range transactions {
		items = append(items, dto.VencimentoItemDTO{
			ID:          t.ID.String(),
			Type:        "transaction",
			Description: t.Description,
			Amount:      t.Amount,
			Date:        t.Date.Format("2006-01-02"),
			Status:      t.Status,
		})
	}

	// Buscar faturas fechadas não pagas (status 'closed')
	creditCards, err := s.creditCardRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, cc := range creditCards {
		bills, err := s.creditCardBillRepo.FindByCreditCardID(ctx, cc.ID)
		if err != nil {
			continue
		}

		for _, bill := range bills {
			if bill.Status == "closed" && bill.DueDate.Before(endDate) && bill.DueDate.After(time.Now().AddDate(0, 0, -1)) {
				total, err := s.creditCardBillRepo.CalculateBillTotal(ctx, bill.ID)
				if err != nil {
					continue
				}

				items = append(items, dto.VencimentoItemDTO{
					ID:          bill.ID.String(),
					Type:        "bill",
					Description: fmt.Sprintf("Fatura %s - %02d/%d", cc.Name, bill.Month, bill.Year),
					Amount:      total,
					Date:        bill.DueDate.Format("2006-01-02"),
					Status:      bill.Status,
				})
			}
		}
	}

	return &dto.CalendarioVencimentosResponse{
		Items: items,
	}, nil
}

// GetGastosPorTag agrupa gastos por tag usando unnest do Postgres
func (s *analyticsService) GetGastosPorTag(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*dto.GastosPorTagResponse, error) {
	filters := &dto.TransactionFilters{
		Type:      func() *string { t := "expense"; return &t }(),
		StartDate: func() *string { s := startDate.Format("2006-01-02"); return &s }(),
		EndDate:   func() *string { s := endDate.Format("2006-01-02"); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Agrupar por tag
	tagMap := make(map[string]*dto.TagBreakdownDTO)
	var totalAmount int64

	for _, t := range transactions {
		if t.Tags == nil || len(t.Tags) == 0 {
			continue
		}

		for _, tag := range t.Tags {
			if breakdown, exists := tagMap[tag]; exists {
				breakdown.TotalAmount += t.Amount
				breakdown.TransactionCount++
			} else {
				tagMap[tag] = &dto.TagBreakdownDTO{
					Tag:              tag,
					TotalAmount:      t.Amount,
					TransactionCount: 1,
				}
			}
			totalAmount += t.Amount
		}
	}

	// Converter para slice e calcular percentuais
	breakdowns := make([]dto.TagBreakdownDTO, 0, len(tagMap))
	for _, breakdown := range tagMap {
		if totalAmount > 0 {
			breakdown.Percentage = float64(breakdown.TotalAmount) / float64(totalAmount) * 100
		}
		breakdowns = append(breakdowns, *breakdown)
	}

	return &dto.GastosPorTagResponse{
		Tags: breakdowns,
	}, nil
}
