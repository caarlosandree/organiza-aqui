package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type CreditCardService interface {
	CreateCreditCard(ctx context.Context, userID uuid.UUID, req *dto.CreateCreditCardRequest) (*dto.CreditCardDTO, error)
	GetCreditCardByID(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) (*dto.CreditCardDTO, error)
	GetCreditCardsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.CreditCardDTO, error)
	UpdateCreditCard(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID, req *dto.UpdateCreditCardRequest) (*dto.CreditCardDTO, error)
	DeleteCreditCard(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) error
	GetAvailableLimit(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) (int64, error)
	GetBillProjection(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) (*dto.ProjecaoFaturaResponse, error)
	GetBills(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) ([]*dto.CreditCardBillDTO, error)
	CloseBill(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID, billID uuid.UUID) error
	PayBill(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID, billID uuid.UUID, req *dto.PayBillRequest) (*dto.CreditCardBillDTO, error)
}

type creditCardService struct {
	creditCardRepo     repository.CreditCardRepository
	creditCardBillRepo repository.CreditCardBillRepository
	accountRepo        repository.AccountRepository
	transactionRepo    repository.TransactionRepository
}

func NewCreditCardService(
	creditCardRepo repository.CreditCardRepository,
	creditCardBillRepo repository.CreditCardBillRepository,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
) CreditCardService {
	return &creditCardService{
		creditCardRepo:     creditCardRepo,
		creditCardBillRepo: creditCardBillRepo,
		accountRepo:        accountRepo,
		transactionRepo:    transactionRepo,
	}
}

func (s *creditCardService) CreateCreditCard(ctx context.Context, userID uuid.UUID, req *dto.CreateCreditCardRequest) (*dto.CreditCardDTO, error) {
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

	now := time.Now()
	creditCard := &model.CreditCard{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		AccountID:   accountID,
		LimitAmount: req.LimitAmount,
		ClosingDay:  req.ClosingDay,
		DueDay:      req.DueDay,
		Color:       req.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.creditCardRepo.Create(ctx, creditCard); err != nil {
		return nil, fmt.Errorf("erro ao criar cartão: %w", err)
	}

	return s.modelToDTO(creditCard), nil
}

func (s *creditCardService) GetCreditCardByID(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) (*dto.CreditCardDTO, error) {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return nil, fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(creditCard), nil
}

func (s *creditCardService) GetCreditCardsByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.CreditCardDTO, error) {
	creditCards, err := s.creditCardRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cartões: %w", err)
	}

	dtos := make([]*dto.CreditCardDTO, len(creditCards))
	for i, cc := range creditCards {
		dtos[i] = s.modelToDTO(cc)
	}

	return dtos, nil
}

func (s *creditCardService) UpdateCreditCard(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID, req *dto.UpdateCreditCardRequest) (*dto.CreditCardDTO, error) {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return nil, fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

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

	creditCard.Name = req.Name
	creditCard.AccountID = accountID
	creditCard.LimitAmount = req.LimitAmount
	creditCard.ClosingDay = req.ClosingDay
	creditCard.DueDay = req.DueDay
	creditCard.Color = req.Color
	creditCard.UpdatedAt = time.Now()

	if err := s.creditCardRepo.Update(ctx, creditCard); err != nil {
		return nil, fmt.Errorf("erro ao atualizar cartão: %w", err)
	}

	return s.modelToDTO(creditCard), nil
}

func (s *creditCardService) DeleteCreditCard(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) error {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	if err := s.creditCardRepo.Delete(ctx, creditCardID); err != nil {
		return fmt.Errorf("erro ao deletar cartão: %w", err)
	}

	return nil
}

// GetAvailableLimit calcula o limite disponível do cartão
func (s *creditCardService) GetAvailableLimit(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) (int64, error) {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return 0, fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return 0, appError.ErrUnauthorizedAccess
	}

	usedLimit, err := s.creditCardRepo.CalculateUsedLimit(ctx, creditCardID)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular limite utilizado: %w", err)
	}

	return creditCard.LimitAmount - usedLimit, nil
}

// GetBillProjection projeta a fatura futura baseada em transações pendentes
func (s *creditCardService) GetBillProjection(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) (*dto.ProjecaoFaturaResponse, error) {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return nil, fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	// Calcular total de transações pendentes na conta do cartão
	usedLimit, err := s.creditCardRepo.CalculateUsedLimit(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular limite utilizado: %w", err)
	}

	// Calcular datas de fechamento e vencimento do próximo mês
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// Se já passou o dia de fechamento, projetar para o próximo mês
	if now.Day() >= creditCard.ClosingDay {
		month++
		if month > 12 {
			month = 1
			year++
		}
	}

	closingDate := time.Date(year, time.Month(month), creditCard.ClosingDay, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(year, time.Month(month), creditCard.DueDay, 0, 0, 0, 0, time.UTC)

	return &dto.ProjecaoFaturaResponse{
		CreditCardID:   creditCardID.String(),
		Month:          month,
		Year:           year,
		ProjectedAmount: usedLimit,
		ClosingDate:    closingDate.Format("2006-01-02"),
		DueDate:        dueDate.Format("2006-01-02"),
	}, nil
}

// GetBills retorna todas as faturas de um cartão
func (s *creditCardService) GetBills(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID) ([]*dto.CreditCardBillDTO, error) {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return nil, fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	bills, err := s.creditCardBillRepo.FindByCreditCardID(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faturas: %w", err)
	}

	dtos := make([]*dto.CreditCardBillDTO, len(bills))
	for i, bill := range bills {
		// Calcular total dinamicamente
		total, err := s.creditCardBillRepo.CalculateBillTotal(ctx, bill.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao calcular total da fatura: %w", err)
		}

		dtos[i] = s.billModelToDTO(bill, total)
	}

	return dtos, nil
}

// CloseBill fecha uma fatura (muda status de 'open' para 'closed')
func (s *creditCardService) CloseBill(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID, billID uuid.UUID) error {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	bill, err := s.creditCardBillRepo.FindByID(ctx, billID)
	if err != nil {
		return fmt.Errorf("erro ao buscar fatura: %w", err)
	}
	if bill == nil {
		return fmt.Errorf("fatura não encontrada: %w", appError.ErrBillNotFound)
	}
	if bill.CreditCardID != creditCardID {
		return appError.ErrUnauthorizedAccess
	}

	if bill.Status != "open" {
		return appError.ErrAlreadyClosed
	}

	if err := s.creditCardBillRepo.CloseBill(ctx, billID); err != nil {
		return fmt.Errorf("erro ao fechar fatura: %w", err)
	}

	return nil
}

// PayBill paga uma fatura criando uma transação de despesa na conta bancária
func (s *creditCardService) PayBill(ctx context.Context, userID uuid.UUID, creditCardID uuid.UUID, billID uuid.UUID, req *dto.PayBillRequest) (*dto.CreditCardBillDTO, error) {
	creditCard, err := s.creditCardRepo.FindByID(ctx, creditCardID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cartão: %w", err)
	}
	if creditCard == nil {
		return nil, fmt.Errorf("cartão %s: %w", creditCardID, appError.ErrCreditCardNotFound)
	}
	if creditCard.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	bill, err := s.creditCardBillRepo.FindByID(ctx, billID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fatura: %w", err)
	}
	if bill == nil {
		return nil, fmt.Errorf("fatura não encontrada: %w", appError.ErrBillNotFound)
	}
	if bill.CreditCardID != creditCardID {
		return nil, appError.ErrUnauthorizedAccess
	}

	if bill.Status != "closed" {
		return nil, appError.ErrNotClosed
	}

	// Validar conta de pagamento
	paymentAccountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account_id inválido: %w", err)
	}

	paymentAccount, err := s.accountRepo.FindByID(ctx, paymentAccountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta de pagamento: %w", err)
	}
	if paymentAccount == nil {
		return nil, errors.New("conta de pagamento não encontrada")
	}
	if paymentAccount.UserID != userID {
		return nil, errors.New("conta de pagamento não pertence ao usuário")
	}

	// Calcular total da fatura
	totalAmount, err := s.creditCardBillRepo.CalculateBillTotal(ctx, billID)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular total da fatura: %w", err)
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	// Criar transação de despesa para pagar a fatura
	now := time.Now()
	paymentTransaction := &model.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		AccountID:   paymentAccountID,
		CategoryID:  nil,
		Type:        "expense",
		Amount:      totalAmount,
		Description: fmt.Sprintf("Pagamento fatura cartão %s - %02d/%d", creditCard.Name, bill.Month, bill.Year),
		Date:        date,
		Status:      "paid",
		CreatedAt:   now,
	}

	if err := s.transactionRepo.Create(ctx, paymentTransaction); err != nil {
		return nil, fmt.Errorf("erro ao criar transação de pagamento: %w", err)
	}

	// Atualizar saldo da conta de pagamento
	if err := s.accountRepo.UpdateBalance(ctx, paymentAccountID, -totalAmount); err != nil {
		return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
	}

	// Atualizar fatura com payment_transaction_id e status 'paid'
	bill.PaymentTransactionID = &paymentTransaction.ID
	bill.Status = "paid"
	bill.UpdatedAt = time.Now()

	if err := s.creditCardBillRepo.Update(ctx, bill); err != nil {
		return nil, fmt.Errorf("erro ao atualizar fatura: %w", err)
	}

	return s.billModelToDTO(bill, totalAmount), nil
}

// modelToDTO converte model.CreditCard para dto.CreditCardDTO
func (s *creditCardService) modelToDTO(creditCard *model.CreditCard) *dto.CreditCardDTO {
	return &dto.CreditCardDTO{
		ID:          creditCard.ID.String(),
		UserID:      creditCard.UserID.String(),
		Name:        creditCard.Name,
		AccountID:   creditCard.AccountID.String(),
		LimitAmount: creditCard.LimitAmount,
		ClosingDay:  creditCard.ClosingDay,
		DueDay:      creditCard.DueDay,
		Color:       creditCard.Color,
		CreatedAt:   creditCard.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   creditCard.UpdatedAt.Format(time.RFC3339),
	}
}

// billModelToDTO converte model.CreditCardBill para dto.CreditCardBillDTO
func (s *creditCardService) billModelToDTO(bill *model.CreditCardBill, totalAmount int64) *dto.CreditCardBillDTO {
	dto := &dto.CreditCardBillDTO{
		ID:          bill.ID.String(),
		CreditCardID: bill.CreditCardID.String(),
		Month:       bill.Month,
		Year:        bill.Year,
		Status:      bill.Status,
		ClosingDate: bill.ClosingDate.Format("2006-01-02"),
		DueDate:     bill.DueDate.Format("2006-01-02"),
		TotalAmount: totalAmount,
		CreatedAt:   bill.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   bill.UpdatedAt.Format(time.RFC3339),
	}

	if bill.PaymentTransactionID != nil {
		paymentIDStr := bill.PaymentTransactionID.String()
		dto.PaymentTransactionID = &paymentIDStr
	}

	return dto
}
