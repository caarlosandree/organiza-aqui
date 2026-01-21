package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type TransactionPeriodService interface {
	GetOrCreatePeriod(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, referenceMonth time.Time, periodType string) (*model.TransactionPeriod, error)
	GetPeriodsByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionPeriodFilters) ([]*dto.TransactionPeriodDTO, error)
	GetPeriodWithTransactions(ctx context.Context, userID uuid.UUID, periodID uuid.UUID) (*dto.PeriodWithTransactionsDTO, error)
}

type transactionPeriodService struct {
	periodRepo      repository.TransactionPeriodRepository
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewTransactionPeriodService(
	periodRepo repository.TransactionPeriodRepository,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
) TransactionPeriodService {
	return &transactionPeriodService{
		periodRepo:      periodRepo,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// GetOrCreatePeriod busca um período existente ou cria um novo baseado em account_id, period_type, year, month
func (s *transactionPeriodService) GetOrCreatePeriod(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, referenceMonth time.Time, periodType string) (*model.TransactionPeriod, error) {
	// Garantir que referenceMonth seja o primeiro dia do mês
	referenceMonth = time.Date(referenceMonth.Year(), referenceMonth.Month(), 1, 0, 0, 0, 0, referenceMonth.Location())

	year := referenceMonth.Year()
	month := int(referenceMonth.Month())

	// Buscar período existente
	period, err := s.periodRepo.FindByAccountAndPeriod(ctx, accountID, periodType, year, month)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("erro ao buscar período: %w", err)
	}

	if period != nil {
		return period, nil
	}

	// Criar novo período
	now := time.Now()
	newPeriod := &model.TransactionPeriod{
		ID:         uuid.New(),
		UserID:     userID,
		AccountID:  accountID,
		PeriodType: periodType,
		Year:       year,
		Month:      month,
		Status:     "open",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.periodRepo.Create(ctx, newPeriod); err != nil {
		return nil, fmt.Errorf("erro ao criar período: %w", err)
	}

	return newPeriod, nil
}

// GetPeriodsByUserID lista períodos do usuário com filtros opcionais
func (s *transactionPeriodService) GetPeriodsByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionPeriodFilters) ([]*dto.TransactionPeriodDTO, error) {
	periods, err := s.periodRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar períodos: %w", err)
	}

	dtos := make([]*dto.TransactionPeriodDTO, len(periods))
	for i, period := range periods {
		dtos[i] = s.modelToDTO(period)

		// Buscar nome da conta
		account, err := s.accountRepo.FindByID(ctx, period.AccountID)
		if err == nil && account != nil {
			accountName := account.Name
			dtos[i].AccountName = &accountName
		}
	}

	return dtos, nil
}

// GetPeriodWithTransactions obtém um período com suas transações e estatísticas
func (s *transactionPeriodService) GetPeriodWithTransactions(ctx context.Context, userID uuid.UUID, periodID uuid.UUID) (*dto.PeriodWithTransactionsDTO, error) {
	// Buscar período
	period, err := s.periodRepo.FindByID(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar período: %w", err)
	}
	if period == nil {
		return nil, fmt.Errorf("período %s: %w", periodID, appError.ErrNotFound)
	}
	if period.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Buscar transações do período
	transactionFilters := &dto.TransactionFilters{
		AccountID: func() *string { s := period.AccountID.String(); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, transactionFilters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transações: %w", err)
	}

	// Filtrar transações do período
	var periodTransactions []*model.Transaction
	for _, t := range transactions {
		if t.PeriodID != nil && *t.PeriodID == periodID {
			periodTransactions = append(periodTransactions, t)
		}
	}

	// Calcular estatísticas
	var totalIncome, totalBankExpense, totalCreditCardExpense int64

	for _, t := range periodTransactions {
		if period.PeriodType == "bank" {
			if t.Type == "income" {
				totalIncome += t.Amount
			} else if t.Type == "expense" {
				totalBankExpense += t.Amount
			}
		} else if period.PeriodType == "credit_card" {
			if t.Type == "expense" {
				totalCreditCardExpense += t.Amount
			}
		}
	}

	balance := totalIncome - totalBankExpense

	// Converter transações para DTOs
	transactionDTOs := make([]*dto.TransactionDTO, len(periodTransactions))
	for i, t := range periodTransactions {
		transactionDTOs[i] = s.transactionToDTO(t)
	}

	// Buscar nome da conta
	account, err := s.accountRepo.FindByID(ctx, period.AccountID)
	accountName := ""
	if err == nil && account != nil {
		accountName = account.Name
	}

	return &dto.PeriodWithTransactionsDTO{
		Period:                 *s.modelToDTOWithAccountName(period, accountName),
		Transactions:           transactionDTOs,
		TotalIncome:            totalIncome,
		TotalBankExpense:       totalBankExpense,
		TotalCreditCardExpense: totalCreditCardExpense,
		Balance:                balance,
	}, nil
}

func (s *transactionPeriodService) modelToDTO(period *model.TransactionPeriod) *dto.TransactionPeriodDTO {
	return &dto.TransactionPeriodDTO{
		ID:         period.ID.String(),
		UserID:     period.UserID.String(),
		AccountID:  period.AccountID.String(),
		PeriodType: period.PeriodType,
		Year:       period.Year,
		Month:      period.Month,
		Status:     period.Status,
		CreatedAt:  period.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  period.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *transactionPeriodService) modelToDTOWithAccountName(period *model.TransactionPeriod, accountName string) *dto.TransactionPeriodDTO {
	dto := s.modelToDTO(period)
	dto.AccountName = &accountName
	return dto
}

func (s *transactionPeriodService) transactionToDTO(transaction *model.Transaction) *dto.TransactionDTO {
	dto := &dto.TransactionDTO{
		ID:          transaction.ID.String(),
		UserID:      transaction.UserID.String(),
		AccountID:   transaction.AccountID.String(),
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date.Format("2006-01-02"),
		Status:      transaction.Status,
		CreatedAt:   transaction.CreatedAt.Format(time.RFC3339),
	}

	if transaction.CategoryID != nil {
		categoryIDStr := transaction.CategoryID.String()
		dto.CategoryID = &categoryIDStr
	}

	if len(transaction.Tags) > 0 {
		dto.Tags = []string(transaction.Tags)
	}

	if transaction.ToAccountID != nil {
		toAccountIDStr := transaction.ToAccountID.String()
		dto.ToAccountID = &toAccountIDStr
	}

	if transaction.ParentTransactionID != nil {
		parentIDStr := transaction.ParentTransactionID.String()
		dto.ParentTransactionID = &parentIDStr
	}

	if transaction.InstallmentNumber != nil {
		dto.InstallmentNumber = transaction.InstallmentNumber
	}

	if transaction.TotalInstallments != nil {
		dto.TotalInstallments = transaction.TotalInstallments
	}

	if transaction.PeriodID != nil {
		periodIDStr := transaction.PeriodID.String()
		dto.PeriodID = &periodIDStr
	}

	if transaction.ReferenceMonth != nil {
		referenceMonthStr := transaction.ReferenceMonth.Format("2006-01")
		dto.ReferenceMonth = &referenceMonthStr
	}

	return dto
}
