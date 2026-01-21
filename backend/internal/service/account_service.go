package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type AccountService interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, req *dto.CreateAccountRequest) (*dto.AccountDTO, error)
	GetAccountByID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*dto.AccountDTO, error)
	GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.AccountDTO, error)
	UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req *dto.UpdateAccountRequest) (*dto.AccountDTO, error)
	DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
	UpdateInitialBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req *dto.UpdateInitialBalanceRequest) (*dto.AccountDTO, error)
	RecalculateBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*dto.AccountDTO, error)
}

type accountService struct {
	db              *sqlx.DB
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewAccountService(db *sqlx.DB, accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) AccountService {
	return &accountService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, userID uuid.UUID, req *dto.CreateAccountRequest) (*dto.AccountDTO, error) {
	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		return nil, fmt.Errorf("bank_id inválido: %w", err)
	}

	now := time.Now()
	account := &model.Account{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Balance:   0,
		Currency:  req.Currency,
		BankID:    bankID,
		CreatedAt: now,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("erro ao criar conta: %w", err)
	}

	return s.modelToDTO(account), nil
}

func (s *accountService) GetAccountByID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*dto.AccountDTO, error) {
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

	return s.modelToDTO(account), nil
}

func (s *accountService) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.AccountDTO, error) {
	accounts, err := s.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contas: %w", err)
	}

	dtos := make([]*dto.AccountDTO, len(accounts))
	for i, account := range accounts {
		dtos[i] = s.modelToDTO(account)
	}

	return dtos, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req *dto.UpdateAccountRequest) (*dto.AccountDTO, error) {
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

	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		return nil, fmt.Errorf("bank_id inválido: %w", err)
	}

	account.Name = req.Name
	account.Type = req.Type
	account.Currency = req.Currency
	account.BankID = bankID

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("erro ao atualizar conta: %w", err)
	}

	return s.modelToDTO(account), nil
}

func (s *accountService) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("erro ao buscar conta: %w", err)
	}
	if account == nil {
		return fmt.Errorf("conta %s: %w", accountID, appError.ErrAccountNotFound)
	}
	if account.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.accountRepo.Delete(ctx, accountID); err != nil {
		return fmt.Errorf("erro ao deletar conta: %w", err)
	}

	return nil
}

func (s *accountService) UpdateInitialBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req *dto.UpdateInitialBalanceRequest) (*dto.AccountDTO, error) {
	// Validar conta e permissões
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

	// Parse da data
	referenceDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	// Validar que a data não é futura
	now := time.Now()
	if referenceDate.After(now) {
		return nil, appError.ErrReferenceDateFuture
	}

	// CRÍTICO: Para cartão de crédito, o usuário digita um valor positivo (dívida)
	// mas o sistema armazena como negativo. Inverter o sinal antes de processar.
	balanceToStore := req.Balance
	if account.Type == "credit" && balanceToStore > 0 {
		// Usuário digitou "1500" (pensando em dívida de R$ 1500)
		// Sistema armazena como -1500 (saldo negativo = dívida)
		balanceToStore = -balanceToStore
	}

	// Calcular saldo de transações até a data de referência
	// GetAccountBalanceAtDate retorna a soma de transações (income - expense) até aquela data
	transactionsAtReferenceDate, err := s.transactionRepo.GetAccountBalanceAtDate(ctx, accountID, referenceDate)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular saldo de transações na data de referência: %w", err)
	}

	// Calcular saldo esperado na data de referência
	// Se já existe initial_balance, considerar ele; senão, considerar 0
	var currentInitialBalance int64
	if account.InitialBalance != nil {
		currentInitialBalance = *account.InitialBalance
	}

	// Saldo esperado = initial_balance atual + transações até a data de referência
	expectedBalanceAtDate := currentInitialBalance + transactionsAtReferenceDate

	// Calcular diferença entre o saldo informado e o saldo esperado
	difference := balanceToStore - expectedBalanceAtDate

	// Usar transação ACID para garantir consistência
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Se houver diferença, criar transação de ajuste
	if difference != 0 {
		adjustmentTransaction := &model.Transaction{
			ID:          uuid.New(),
			UserID:      userID,
			AccountID:   accountID,
			CategoryID:  nil,
			Type:        "adjustment",
			Amount:      difference,
			Description: "Ajuste de saldo inicial",
			Date:        referenceDate,
			Status:      "paid",
			Tags:        nil,
			CreatedAt:   time.Now(),
		}

		// Criar transação de ajuste
		createQuery := `
			INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date,
			                         status, tags, created_at)
			VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
			        :status, :tags, :created_at)
		`
		_, err = tx.NamedExecContext(ctx, createQuery, adjustmentTransaction)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar transação de ajuste: %w", err)
		}

		// Atualizar saldo da conta
		updateBalanceQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, updateBalanceQuery, difference, accountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
		}
	}

	// Atualizar saldo inicial e data de referência (usar balanceToStore, não req.Balance)
	updateInitialBalanceQuery := `UPDATE accounts SET initial_balance = $1, initial_balance_date = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateInitialBalanceQuery, balanceToStore, referenceDate, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar saldo inicial: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	// Buscar conta atualizada
	updatedAccount, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta atualizada: %w", err)
	}

	return s.modelToDTO(updatedAccount), nil
}

func (s *accountService) RecalculateBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*dto.AccountDTO, error) {
	// Validar conta e permissões
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

	// Calcular saldo baseado no saldo inicial e nas transações
	calculatedBalance, err := s.transactionRepo.CalculateAccountBalanceFromTransactions(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular saldo: %w", err)
	}

	// Atualizar saldo da conta
	if err := s.accountRepo.SetBalance(ctx, accountID, calculatedBalance); err != nil {
		return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
	}

	// Buscar conta atualizada
	updatedAccount, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta atualizada: %w", err)
	}

	return s.modelToDTO(updatedAccount), nil
}

func (s *accountService) modelToDTO(account *model.Account) *dto.AccountDTO {
	dto := &dto.AccountDTO{
		ID:        account.ID.String(),
		UserID:    account.UserID.String(),
		Name:      account.Name,
		Type:      account.Type,
		Balance:   account.Balance,
		Currency:  account.Currency,
		BankID:    account.BankID.String(),
		CreatedAt: account.CreatedAt.Format(time.RFC3339),
	}

	if account.InitialBalance != nil {
		dto.InitialBalance = account.InitialBalance
	}

	if account.InitialBalanceDate != nil {
		dto.InitialBalanceDate = func() *string {
			s := account.InitialBalanceDate.Format(time.RFC3339)
			return &s
		}()
	}

	return dto
}
