package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, userID uuid.UUID, req *dto.CreateTransactionRequest) (*dto.TransactionDTO, error)
	GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*dto.TransactionDTO, error)
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionFilters) ([]*dto.TransactionDTO, int, error)
	UpdateTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionDTO, error)
	DeleteTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error
	GetStatement(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, startDate, endDate time.Time) (*dto.StatementResponse, error)
}

type transactionService struct {
	db                *sqlx.DB
	transactionRepo   repository.TransactionRepository
	accountRepo       repository.AccountRepository
	categoryRepo      repository.CategoryRepository
	periodService     TransactionPeriodService
	installmentService InstallmentService
}

func NewTransactionService(
	db *sqlx.DB,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	periodService TransactionPeriodService,
	installmentService InstallmentService,
) TransactionService {
	return &transactionService{
		db:                db,
		transactionRepo:   transactionRepo,
		accountRepo:       accountRepo,
		categoryRepo:      categoryRepo,
		periodService:     periodService,
		installmentService: installmentService,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, userID uuid.UUID, req *dto.CreateTransactionRequest) (*dto.TransactionDTO, error) {
	// Validar conta
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account_id inválido: %w", err)
	}

	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("conta %s: %w", accountID, appError.ErrAccountNotFound)
	}
	if account.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Validar categoria se fornecida
	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		parsedCategoryID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("category_id inválido: %w", err)
		}

		category, err := s.categoryRepo.FindByID(ctx, parsedCategoryID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
		}
		if category == nil {
			return nil, fmt.Errorf("categoria %s: %w", parsedCategoryID, appError.ErrCategoryNotFound)
		}
		if category.UserID != userID {
			return nil, appError.ErrUnauthorizedAccess
		}

		categoryID = &parsedCategoryID
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	// Determinar mês de referência
	var referenceMonth time.Time
	if req.ReferenceMonth != nil {
		// Parse do mês de referência (formato YYYY-MM)
		referenceMonth, err = time.Parse("2006-01", *req.ReferenceMonth)
		if err != nil {
			return nil, fmt.Errorf("reference_month inválido: %w", err)
		}
		// Garantir que seja o primeiro dia do mês
		referenceMonth = time.Date(referenceMonth.Year(), referenceMonth.Month(), 1, 0, 0, 0, 0, referenceMonth.Location())
	} else {
		// Se não fornecido, usar o mês da data da transação
		referenceMonth = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	}

	// Determinar tipo de período baseado no tipo da conta
	periodType := "bank"
	if account.Type == "credit" {
		periodType = "credit_card"
	}

	// Obter ou criar período
	period, err := s.periodService.GetOrCreatePeriod(ctx, userID, accountID, referenceMonth, periodType)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter/criar período: %w", err)
	}

	// Definir status padrão
	status := "paid"
	if req.Status != "" {
		status = req.Status
	}

	// Processar tags
	var tags pq.StringArray
	if req.Tags != nil {
		tags = pq.StringArray(req.Tags)
	} else {
		tags = pq.StringArray{}
	}

	now := time.Now()

	// Se for transferência, criar duas transações
	if req.Type == "transfer" {
		if req.ToAccountID == nil {
			return nil, errors.New("to_account_id é obrigatório para transferências")
		}

		toAccountID, err := uuid.Parse(*req.ToAccountID)
		if err != nil {
			return nil, fmt.Errorf("to_account_id inválido: %w", err)
		}

		// Validar conta destino
		toAccount, err := s.accountRepo.FindByID(ctx, toAccountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar conta destino: %w", err)
		}
		if toAccount == nil {
			return nil, fmt.Errorf("conta destino %s: %w", toAccountID, appError.ErrAccountNotFound)
		}
		if toAccount.UserID != userID {
			return nil, appError.ErrUnauthorizedAccess
		}

		// Validar que não é a mesma conta
		if accountID == toAccountID {
			return nil, errors.New("conta origem e destino não podem ser a mesma")
		}

		// Determinar períodos para ambas as contas
		fromPeriodType := "bank"
		if account.Type == "credit" {
			fromPeriodType = "credit_card"
		}
		fromPeriod, err := s.periodService.GetOrCreatePeriod(ctx, userID, accountID, referenceMonth, fromPeriodType)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter/criar período origem: %w", err)
		}

		toPeriodType := "bank"
		if toAccount.Type == "credit" {
			toPeriodType = "credit_card"
		}
		toPeriod, err := s.periodService.GetOrCreatePeriod(ctx, userID, toAccountID, referenceMonth, toPeriodType)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter/criar período destino: %w", err)
		}

		// Criar transação de saída
		fromTransaction := &model.Transaction{
			ID:             uuid.New(),
			UserID:         userID,
			AccountID:      accountID,
			CategoryID:     categoryID,
			Type:           "transfer",
			Amount:         req.Amount,
			Description:    req.Description,
			Date:           date,
			Status:         status,
			Tags:           tags,
			ToAccountID:    &toAccountID,
			PeriodID:       &fromPeriod.ID,
			ReferenceMonth: &referenceMonth,
			CreatedAt:      now,
		}

		// Criar transação de entrada
		toTransaction := &model.Transaction{
			ID:             uuid.New(),
			UserID:         userID,
			AccountID:      toAccountID,
			CategoryID:     categoryID,
			Type:           "transfer",
			Amount:         req.Amount,
			Description:    req.Description,
			Date:           date,
			Status:         status,
			Tags:           tags,
			ToAccountID:    &accountID,
			PeriodID:       &toPeriod.ID,
			ReferenceMonth: &referenceMonth,
			CreatedAt:      now,
		}

		// Criar transferência usando transação ACID
		if err := s.transactionRepo.CreateTransfer(ctx, fromTransaction, toTransaction, accountID, toAccountID, req.Amount); err != nil {
			return nil, fmt.Errorf("erro ao criar transferência: %w", err)
		}

		return s.modelToDTO(fromTransaction), nil
	}

	// Para income, expense e adjustment, criar transação normal
	var balanceAdjustment int64
	if req.Type == "income" {
		balanceAdjustment = req.Amount
	} else if req.Type == "expense" {
		balanceAdjustment = -req.Amount
	} else if req.Type == "adjustment" {
		// Transações de ajuste: aplicar o valor diretamente
		balanceAdjustment = req.Amount
	}

	// Criar transação
	transaction := &model.Transaction{
		ID:             uuid.New(),
		UserID:         userID,
		AccountID:      accountID,
		CategoryID:     categoryID,
		Type:           req.Type,
		Amount:         req.Amount,
		Description:    req.Description,
		Date:           date,
		Status:         status,
		Tags:           tags,
		PeriodID:       &period.ID,
		ReferenceMonth: &referenceMonth,
		CreatedAt:      now,
	}

	// Se for parcelamento, usar InstallmentService para criar transação mãe e parcelas
	if req.TotalInstallments != nil && *req.TotalInstallments > 1 {
		return s.installmentService.CreateInstallments(ctx, userID, req)
	}

	// Usar transação ACID para garantir consistência
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Criar transação
	query := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date, 
		                         status, tags, to_account_id, parent_transaction_id, installment_number, 
		                         total_installments, external_id, period_id, reference_month, created_at)
		VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
		        :status, :tags, :to_account_id, :parent_transaction_id, :installment_number,
		        :total_installments, :external_id, :period_id, :reference_month, :created_at)
	`
	_, err = tx.NamedExecContext(ctx, query, transaction)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar transação: %w", err)
	}

	// Atualizar saldo da conta
	if balanceAdjustment != 0 {
		updateQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateQuery, balanceAdjustment, accountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return s.modelToDTO(transaction), nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*dto.TransactionDTO, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação: %w", err)
	}
	if transaction == nil {
		return nil, fmt.Errorf("transação %s: %w", transactionID, appError.ErrTransactionNotFound)
	}
	if transaction.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(transaction), nil
}

func (s *transactionService) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, filters *dto.TransactionFilters) ([]*dto.TransactionDTO, int, error) {
	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar transações: %w", err)
	}

	total, err := s.transactionRepo.CountByUserID(ctx, userID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar transações: %w", err)
	}

	dtos := make([]*dto.TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		dtos[i] = s.modelToDTO(transaction)
	}

	return dtos, total, nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionDTO, error) {
	// Buscar transação existente
	oldTransaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação: %w", err)
	}
	if oldTransaction == nil {
		return nil, errors.New("transação não encontrada")
	}
	if oldTransaction.UserID != userID {
		return nil, errors.New("transação não pertence ao usuário")
	}

	// Validar nova conta
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account_id inválido: %w", err)
	}

	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("conta %s: %w", accountID, appError.ErrAccountNotFound)
	}
	if account.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Verificar se o valor, tipo ou conta mudou - se não mudou, não precisa atualizar saldo
	valueOrTypeChanged := oldTransaction.Amount != req.Amount || oldTransaction.Type != req.Type || oldTransaction.AccountID != accountID

	// Calcular ajustes de saldo (será feito na transação)
	var oldBalanceAdjustment int64
	if valueOrTypeChanged {
		if oldTransaction.Type == "income" {
			oldBalanceAdjustment = -oldTransaction.Amount
		} else if oldTransaction.Type == "expense" {
			oldBalanceAdjustment = oldTransaction.Amount
		} else if oldTransaction.Type == "adjustment" {
			// Transações de ajuste: reverter o valor (subtrair o que foi adicionado)
			oldBalanceAdjustment = -oldTransaction.Amount
		}
	}

	// Validar categoria se fornecida
	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		parsedCategoryID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("category_id inválido: %w", err)
		}

		category, err := s.categoryRepo.FindByID(ctx, parsedCategoryID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
		}
		if category == nil {
			return nil, fmt.Errorf("categoria %s: %w", parsedCategoryID, appError.ErrCategoryNotFound)
		}
		if category.UserID != userID {
			return nil, appError.ErrUnauthorizedAccess
		}

		categoryID = &parsedCategoryID
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	// Determinar mês de referência
	var referenceMonth time.Time
	if req.ReferenceMonth != nil {
		// Parse do mês de referência (formato YYYY-MM)
		referenceMonth, err = time.Parse("2006-01", *req.ReferenceMonth)
		if err != nil {
			return nil, fmt.Errorf("reference_month inválido: %w", err)
		}
		// Garantir que seja o primeiro dia do mês
		referenceMonth = time.Date(referenceMonth.Year(), referenceMonth.Month(), 1, 0, 0, 0, 0, referenceMonth.Location())
	} else {
		// Se não fornecido, usar o mês da data da transação
		referenceMonth = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	}

	// Verificar se período precisa ser atualizado (mudou conta, data ou reference_month)
	needsPeriodUpdate := oldTransaction.AccountID != accountID ||
		oldTransaction.Date.Year() != date.Year() || oldTransaction.Date.Month() != date.Month() ||
		(oldTransaction.ReferenceMonth == nil && referenceMonth != time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())) ||
		(oldTransaction.ReferenceMonth != nil && (oldTransaction.ReferenceMonth.Year() != referenceMonth.Year() || oldTransaction.ReferenceMonth.Month() != referenceMonth.Month()))

	var periodID *uuid.UUID
	if needsPeriodUpdate {
		// Determinar tipo de período baseado no tipo da conta
		periodType := "bank"
		if account.Type == "credit" {
			periodType = "credit_card"
		}

		// Obter ou criar novo período
		period, err := s.periodService.GetOrCreatePeriod(ctx, userID, accountID, referenceMonth, periodType)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter/criar período: %w", err)
		}
		periodID = &period.ID
	} else {
		// Manter período existente
		periodID = oldTransaction.PeriodID
		if oldTransaction.ReferenceMonth != nil {
			referenceMonth = *oldTransaction.ReferenceMonth
		} else {
			referenceMonth = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		}
	}

	// Calcular novo ajuste de saldo
	var newBalanceAdjustment int64
	if req.Type == "income" {
		newBalanceAdjustment = req.Amount
	} else if req.Type == "expense" {
		newBalanceAdjustment = -req.Amount
	} else if req.Type == "adjustment" {
		// Transações de ajuste: aplicar o valor diretamente
		newBalanceAdjustment = req.Amount
	}

	// Processar tags
	var tags pq.StringArray
	if req.Tags != nil {
		tags = pq.StringArray(req.Tags)
	} else {
		tags = oldTransaction.Tags
	}

	// Atualizar transação
	oldTransaction.AccountID = accountID
	oldTransaction.CategoryID = categoryID
	oldTransaction.Type = req.Type
	oldTransaction.Amount = req.Amount
	oldTransaction.Description = req.Description
	oldTransaction.Date = date
	oldTransaction.Tags = tags
	oldTransaction.PeriodID = periodID
	oldTransaction.ReferenceMonth = &referenceMonth
	if req.Status != "" {
		oldTransaction.Status = req.Status
	}
	if req.ToAccountID != nil {
		toAccountID, err := uuid.Parse(*req.ToAccountID)
		if err == nil {
			oldTransaction.ToAccountID = &toAccountID
		}
	}

	// Usar transação ACID para garantir consistência
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Atualizar transação
	updateQuery := `
		UPDATE transactions 
		SET account_id = :account_id, category_id = :category_id, type = :type, 
		    amount = :amount, description = :description, date = :date,
		    status = :status, tags = :tags, to_account_id = :to_account_id,
		    parent_transaction_id = :parent_transaction_id, installment_number = :installment_number,
		    total_installments = :total_installments, external_id = :external_id,
		    period_id = :period_id, reference_month = :reference_month
		WHERE id = :id
	`
	_, err = tx.NamedExecContext(ctx, updateQuery, oldTransaction)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar transação: %w", err)
	}

	// Reverter saldo da conta antiga (apenas se valor/tipo/conta mudou)
	if valueOrTypeChanged && oldBalanceAdjustment != 0 {
		revertQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, revertQuery, oldBalanceAdjustment, oldTransaction.AccountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao reverter saldo da conta: %w", err)
		}
	}

	// Aplicar novo saldo (apenas se valor/tipo mudou)
	if valueOrTypeChanged && newBalanceAdjustment != 0 {
		updateBalanceQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateBalanceQuery, newBalanceAdjustment, accountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return s.modelToDTO(oldTransaction), nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error {
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("erro ao buscar transação: %w", err)
	}
	if transaction == nil {
		return fmt.Errorf("transação %s: %w", transactionID, appError.ErrTransactionNotFound)
	}
	if transaction.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	// Calcular ajuste de saldo a ser revertido
	var balanceAdjustment int64
	if transaction.Type == "income" {
		balanceAdjustment = -transaction.Amount
	} else if transaction.Type == "expense" {
		balanceAdjustment = transaction.Amount
	} else if transaction.Type == "adjustment" {
		// Transações de ajuste: reverter o valor (subtrair o que foi adicionado)
		balanceAdjustment = -transaction.Amount
	}

	// Usar transação ACID para garantir consistência
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Reverter saldo da conta
	if balanceAdjustment != 0 {
		revertQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, revertQuery, balanceAdjustment, transaction.AccountID)
		if err != nil {
			return fmt.Errorf("erro ao reverter saldo da conta: %w", err)
		}
	}

	// Deletar transação
	deleteQuery := `DELETE FROM transactions WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, transactionID)
	if err != nil {
		return fmt.Errorf("erro ao deletar transação: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

func (s *transactionService) modelToDTO(transaction *model.Transaction) *dto.TransactionDTO {
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

	if transaction.Tags != nil && len(transaction.Tags) > 0 {
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

func (s *transactionService) GetStatement(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, startDate, endDate time.Time) (*dto.StatementResponse, error) {
	// Verificar se a conta existe e pertence ao usuário
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("conta %s: %w", accountID, appError.ErrAccountNotFound)
	}
	if account.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Buscar transações do período primeiro
	filters := &dto.TransactionFilters{
		AccountID: func() *string { s := accountID.String(); return &s }(),
		StartDate: func() *string { s := startDate.Format("2006-01-02"); return &s }(),
		EndDate:   func() *string { s := endDate.Format("2006-01-02"); return &s }(),
	}

	transactions, err := s.transactionRepo.FindByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transações: %w", err)
	}

	// Buscar saldo inicial (antes do período)
	initialBalance, err := s.transactionRepo.GetAccountBalanceAtDate(ctx, accountID, startDate.AddDate(0, 0, -1))
	if err != nil {
		// Se não conseguir calcular, usar saldo atual da conta e calcular reversamente
		initialBalance = account.Balance

		// Subtrair todas as transações futuras (após o período)
		futureFilters := &dto.TransactionFilters{
			AccountID: func() *string { s := accountID.String(); return &s }(),
			StartDate: func() *string { s := endDate.AddDate(0, 0, 1).Format("2006-01-02"); return &s }(),
		}
		futureTransactions, _ := s.transactionRepo.FindByUserID(ctx, userID, futureFilters)
		for _, t := range futureTransactions {
			if t.Type == "income" {
				initialBalance -= t.Amount
			} else if t.Type == "expense" {
				initialBalance += t.Amount
			}
		}

		// Subtrair também as transações do período
		for _, t := range transactions {
			if t.Type == "income" {
				initialBalance -= t.Amount
			} else if t.Type == "expense" {
				initialBalance += t.Amount
			}
		}
	}

	// Calcular totais
	var totalIncome, totalExpense int64
	var incomeCount, expenseCount int

	for _, t := range transactions {
		if t.Type == "income" {
			totalIncome += t.Amount
			incomeCount++
		} else if t.Type == "expense" {
			totalExpense += t.Amount
			expenseCount++
		}
	}

	// Calcular saldo final
	finalBalance := initialBalance + totalIncome - totalExpense

	// Converter transações para DTOs
	transactionDTOs := make([]*dto.TransactionDTO, len(transactions))
	for i, t := range transactions {
		transactionDTOs[i] = s.modelToDTO(t)
	}

	// Calcular médias
	var avgIncome, avgExpense int64
	if incomeCount > 0 {
		avgIncome = totalIncome / int64(incomeCount)
	}
	if expenseCount > 0 {
		avgExpense = totalExpense / int64(expenseCount)
	}

	return &dto.StatementResponse{
		AccountID:      accountID.String(),
		AccountName:    account.Name,
		StartDate:      startDate.Format("2006-01-02"),
		EndDate:        endDate.Format("2006-01-02"),
		InitialBalance: initialBalance,
		FinalBalance:   finalBalance,
		TotalIncome:    totalIncome,
		TotalExpense:   totalExpense,
		Transactions:   transactionDTOs,
		Summary: dto.StatementSummaryDTO{
			TransactionCount: len(transactions),
			IncomeCount:      incomeCount,
			ExpenseCount:     expenseCount,
			AverageIncome:    avgIncome,
			AverageExpense:   avgExpense,
		},
	}, nil
}
