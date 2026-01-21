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

type InstallmentService interface {
	CreateInstallments(ctx context.Context, userID uuid.UUID, req *dto.CreateTransactionRequest) (*dto.TransactionDTO, error)
	UpdateInstallment(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, req *dto.UpdateInstallmentRequest) (*dto.TransactionDTO, error)
	GetInstallments(ctx context.Context, userID uuid.UUID, parentTransactionID uuid.UUID) ([]*dto.TransactionDTO, error)
	CancelFutureInstallments(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error
}

type installmentService struct {
	db              *sqlx.DB
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
}

func NewInstallmentService(
	db *sqlx.DB,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) InstallmentService {
	return &installmentService{
		db:              db,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

// CreateInstallments cria uma transação parcelada (transação mãe + N parcelas)
func (s *installmentService) CreateInstallments(ctx context.Context, userID uuid.UUID, req *dto.CreateTransactionRequest) (*dto.TransactionDTO, error) {
	if req.TotalInstallments == nil || *req.TotalInstallments <= 1 {
		return nil, errors.New("total_installments deve ser maior que 1 para parcelamento")
	}

	totalInstallments := *req.TotalInstallments

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

	// Calcular valor por parcela
	amountPerInstallment := req.Amount / int64(totalInstallments)
	remainder := req.Amount % int64(totalInstallments) // resto para adicionar na primeira parcela

	// Status padrão
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
	parentID := uuid.New()

	// Iniciar transação ACID
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Criar transação mãe (primeira parcela)
	parentTransaction := &model.Transaction{
		ID:                  parentID,
		UserID:              userID,
		AccountID:           accountID,
		CategoryID:          categoryID,
		Type:                req.Type,
		Amount:              amountPerInstallment + remainder, // primeira parcela recebe o resto
		Description:         req.Description,
		Date:                date,
		Status:              status,
		Tags:                tags,
		ParentTransactionID: nil, // transação mãe não tem pai
		InstallmentNumber:   func() *int { n := 1; return &n }(),
		TotalInstallments:   &totalInstallments,
		CreatedAt:           now,
	}

	// Inserir transação mãe
	queryParent := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date,
		                         status, tags, to_account_id, parent_transaction_id, installment_number,
		                         total_installments, external_id, created_at)
		VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
		        :status, :tags, :to_account_id, :parent_transaction_id, :installment_number,
		        :total_installments, :external_id, :created_at)
	`
	_, err = tx.NamedExecContext(ctx, queryParent, parentTransaction)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar transação mãe: %w", err)
	}

	// Ajustar saldo da conta (apenas primeira parcela)
	if req.Type == "income" {
		queryUpdateBalance := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, queryUpdateBalance, amountPerInstallment+remainder, accountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar saldo: %w", err)
		}
	} else if req.Type == "expense" {
		queryUpdateBalance := `UPDATE accounts SET balance = balance - $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, queryUpdateBalance, amountPerInstallment+remainder, accountID)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar saldo: %w", err)
		}
	}

	// Criar parcelas futuras
	for i := 2; i <= totalInstallments; i++ {
		installmentDate := date.AddDate(0, i-1, 0) // parcelas mensais
		installmentNumber := i

		installment := &model.Transaction{
			ID:                  uuid.New(),
			UserID:              userID,
			AccountID:           accountID,
			CategoryID:          categoryID,
			Type:                req.Type,
			Amount:              amountPerInstallment,
			Description:         fmt.Sprintf("%s (Parcela %d/%d)", req.Description, i, totalInstallments),
			Date:                installmentDate,
			Status:              "pending", // parcelas futuras são pending
			Tags:                tags,
			ParentTransactionID: &parentID,
			InstallmentNumber:   &installmentNumber,
			TotalInstallments:   &totalInstallments,
			CreatedAt:           now,
		}

		queryInstallment := `
			INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, description, date,
			                         status, tags, to_account_id, parent_transaction_id, installment_number,
			                         total_installments, external_id, created_at)
			VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :description, :date,
			        :status, :tags, :to_account_id, :parent_transaction_id, :installment_number,
			        :total_installments, :external_id, :created_at)
		`
		_, err = tx.NamedExecContext(ctx, queryInstallment, installment)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar parcela %d: %w", i, err)
		}
	}

	// Commit transação
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	// Buscar transação mãe criada para retornar
	createdParent, err := s.transactionRepo.FindByID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação criada: %w", err)
	}

	return s.modelToDTO(createdParent), nil
}

// UpdateInstallment atualiza uma parcela com opções de escopo
func (s *installmentService) UpdateInstallment(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, req *dto.UpdateInstallmentRequest) (*dto.TransactionDTO, error) {
	// Buscar transação
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

	if transaction.ParentTransactionID == nil {
		return nil, errors.New("transação não é uma parcela")
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	// Processar tags
	var tags pq.StringArray
	if req.Tags != nil {
		tags = pq.StringArray(req.Tags)
	} else {
		tags = transaction.Tags
	}

	status := transaction.Status
	if req.Status != "" {
		status = req.Status
	}

	// Iniciar transação ACID
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Determinar quais parcelas atualizar baseado no escopo
	var transactionsToUpdate []*model.Transaction

	if req.Scope == "this" {
		// Apenas esta parcela
		transaction.Amount = req.Amount
		transaction.Description = req.Description
		transaction.Date = date
		transaction.Status = status
		transaction.Tags = tags
		transactionsToUpdate = []*model.Transaction{transaction}
	} else if req.Scope == "this_and_future" {
		// Esta e as futuras
		allInstallments, err := s.transactionRepo.FindInstallments(ctx, *transaction.ParentTransactionID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar parcelas: %w", err)
		}

		for _, inst := range allInstallments {
			if inst.InstallmentNumber != nil && *inst.InstallmentNumber >= *transaction.InstallmentNumber {
				inst.Amount = req.Amount
				inst.Description = req.Description
				if inst.InstallmentNumber != nil && *inst.InstallmentNumber == *transaction.InstallmentNumber {
					inst.Date = date
					inst.Status = status
				}
				inst.Tags = tags
				transactionsToUpdate = append(transactionsToUpdate, inst)
			}
		}
	} else if req.Scope == "all" {
		// Todas as parcelas
		allInstallments, err := s.transactionRepo.FindInstallments(ctx, *transaction.ParentTransactionID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar parcelas: %w", err)
		}

		// Buscar também a transação mãe
		parent, err := s.transactionRepo.FindByID(ctx, *transaction.ParentTransactionID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar transação mãe: %w", err)
		}
		if parent != nil {
			parent.Amount = req.Amount
			parent.Description = req.Description
			parent.Date = date
			parent.Status = status
			parent.Tags = tags
			transactionsToUpdate = append(transactionsToUpdate, parent)
		}

		for _, inst := range allInstallments {
			inst.Amount = req.Amount
			inst.Description = req.Description
			if inst.InstallmentNumber != nil && *inst.InstallmentNumber == *transaction.InstallmentNumber {
				inst.Date = date
				inst.Status = status
			}
			inst.Tags = tags
			transactionsToUpdate = append(transactionsToUpdate, inst)
		}
	}

	// Atualizar todas as transações identificadas
	queryUpdate := `
		UPDATE transactions 
		SET amount = :amount, description = :description, date = :date,
		    status = :status, tags = :tags
		WHERE id = :id
	`

	for _, t := range transactionsToUpdate {
		_, err = tx.NamedExecContext(ctx, queryUpdate, t)
		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar parcela: %w", err)
		}
	}

	// Commit transação
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	// Buscar transação atualizada
	updated, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação atualizada: %w", err)
	}

	return s.modelToDTO(updated), nil
}

// GetInstallments retorna todas as parcelas de uma transação mãe
func (s *installmentService) GetInstallments(ctx context.Context, userID uuid.UUID, parentTransactionID uuid.UUID) ([]*dto.TransactionDTO, error) {
	// Verificar se a transação mãe pertence ao usuário
	parent, err := s.transactionRepo.FindByID(ctx, parentTransactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação mãe: %w", err)
	}
	if parent == nil {
		return nil, errors.New("transação mãe não encontrada")
	}
	if parent.UserID != userID {
		return nil, errors.New("transação não pertence ao usuário")
	}

	// Buscar todas as parcelas
	installments, err := s.transactionRepo.FindInstallments(ctx, parentTransactionID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar parcelas: %w", err)
	}

	// Incluir a transação mãe também
	allTransactions := []*model.Transaction{parent}
	allTransactions = append(allTransactions, installments...)

	// Converter para DTOs
	dtos := make([]*dto.TransactionDTO, len(allTransactions))
	for i, t := range allTransactions {
		dtos[i] = s.modelToDTO(t)
	}

	return dtos, nil
}

// CancelFutureInstallments cancela todas as parcelas futuras de uma transação
func (s *installmentService) CancelFutureInstallments(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error {
	// Buscar transação
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

	if transaction.ParentTransactionID == nil {
		return errors.New("transação não é uma parcela")
	}

	// Buscar todas as parcelas
	allInstallments, err := s.transactionRepo.FindInstallments(ctx, *transaction.ParentTransactionID)
	if err != nil {
		return fmt.Errorf("erro ao buscar parcelas: %w", err)
	}

	// Iniciar transação ACID
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// Cancelar parcelas futuras (status = 'cancelled')
	queryCancel := `UPDATE transactions SET status = 'cancelled' WHERE id = $1`
	currentInstallmentNumber := *transaction.InstallmentNumber

	for _, inst := range allInstallments {
		if inst.InstallmentNumber != nil && *inst.InstallmentNumber > currentInstallmentNumber {
			_, err = tx.ExecContext(ctx, queryCancel, inst.ID)
			if err != nil {
				return fmt.Errorf("erro ao cancelar parcela: %w", err)
			}
		}
	}

	// Commit transação
	return tx.Commit()
}

// modelToDTO converte model.Transaction para dto.TransactionDTO
func (s *installmentService) modelToDTO(transaction *model.Transaction) *dto.TransactionDTO {
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

	return dto
}
