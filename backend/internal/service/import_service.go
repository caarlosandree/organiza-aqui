package service

import (
	"context"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
	"github.com/organiza-aqui/backend/internal/util"
)

type ImportService interface {
	ImportOFX(ctx context.Context, userID uuid.UUID, req *dto.ImportOFXRequest) (*dto.ImportResponse, error)
	ImportCSV(ctx context.Context, userID uuid.UUID, req *dto.ImportCSVRequest) (*dto.ImportResponse, error)
	PreviewOFX(ctx context.Context, userID uuid.UUID, req *dto.ImportOFXRequest) (*dto.ImportPreviewResponse, error)
	PreviewCSV(ctx context.Context, userID uuid.UUID, req *dto.ImportCSVRequest) (*dto.ImportPreviewResponse, error)
}

type importService struct {
	db              *sqlx.DB
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	periodService   TransactionPeriodService
}

func NewImportService(
	db *sqlx.DB,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	periodService TransactionPeriodService,
) ImportService {
	return &importService{
		db:              db,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		periodService:   periodService,
	}
}

// generateExternalID gera um hash MD5 para deduplicação
func (s *importService) generateExternalID(date time.Time, amount int64, description string) string {
	// Normalizar descrição: remover espaços extras, converter para minúsculas
	normalizedDesc := strings.ToLower(strings.TrimSpace(description))
	normalizedDesc = strings.Join(strings.Fields(normalizedDesc), " ")

	// Criar string única: data + valor + descrição normalizada
	data := fmt.Sprintf("%s|%d|%s", date.Format("2006-01-02"), amount, normalizedDesc)

	// Gerar hash MD5
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// normalizeDescription normaliza a descrição para comparação
func (s *importService) normalizeDescription(desc string) string {
	normalized := strings.ToLower(strings.TrimSpace(desc))
	return strings.Join(strings.Fields(normalized), " ")
}

// ImportOFX importa transações de um arquivo OFX
func (s *importService) ImportOFX(ctx context.Context, userID uuid.UUID, req *dto.ImportOFXRequest) (*dto.ImportResponse, error) {
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

	// Parse do arquivo OFX
	parser := util.NewOFXParser()
	if err := parser.Parse(req.File); err != nil {
		return nil, fmt.Errorf("erro ao parsear arquivo OFX: %w", err)
	}

	ofxTransactions := parser.GetTransactions()
	if len(ofxTransactions) == 0 {
		return nil, errors.New("nenhuma transação encontrada no arquivo OFX")
	}

	// Obter tipo do arquivo OFX
	ofxFileType := parser.GetFileType()

	// Validar que o tipo da conta corresponde ao tipo do arquivo OFX
	if ofxFileType == "credit_card" && account.Type != "credit" {
		return nil, fmt.Errorf("arquivo OFX é de cartão de crédito, mas a conta selecionada não é do tipo 'credit'")
	}
	if ofxFileType == "bank" && account.Type == "credit" {
		return nil, fmt.Errorf("arquivo OFX é de extrato bancário, mas a conta selecionada é do tipo 'credit'")
	}

	// Iniciar transação ACID
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	var totalProcessed, duplicates, created, errorsCount int
	now := time.Now()
	var earliestDate *time.Time // Data mais antiga do extrato

	// Determinar tipo de período baseado no tipo do arquivo OFX (não apenas da conta)
	// O parser já retorna "credit_card" ou "bank"
	periodType := ofxFileType

	// Processar transações OFX
	for i, ofxTxn := range ofxTransactions {
		totalProcessed++

		// Criar savepoint para isolar cada operação
		savepointName := fmt.Sprintf("sp_%d", i)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %s", savepointName))
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("erro ao criar savepoint: %w", err)
		}

		// Atualizar data mais antiga
		if earliestDate == nil || ofxTxn.Date.Before(*earliestDate) {
			earliestDate = &ofxTxn.Date
		}

		// Usar FITID como external_id para deduplicação (mais confiável que hash)
		externalID := ofxTxn.FITID

		// Verificar se já existe
		existing, err := s.transactionRepo.FindByExternalID(ctx, externalID)
		if err != nil {
			errorsCount++
			continue
		}
		if existing != nil {
			duplicates++
			continue
		}

		// Determinar tipo de transação baseado no TRNTYPE do OFX
		transactionType := "expense"
		if ofxTxn.Type == "CREDIT" {
			transactionType = "income"
		}

		// Usar data da transação como mês de referência (primeiro dia do mês)
		referenceMonth := time.Date(ofxTxn.Date.Year(), ofxTxn.Date.Month(), 1, 0, 0, 0, 0, ofxTxn.Date.Location())

		// Obter ou criar período
		period, err := s.periodService.GetOrCreatePeriod(ctx, userID, accountID, referenceMonth, periodType)
		if err != nil {
			errorsCount++
			continue
		}

		// Criar transação
		transaction := &model.Transaction{
			ID:             uuid.New(),
			UserID:         userID,
			AccountID:      accountID,
			CategoryID:     nil,
			Type:           transactionType,
			Amount:         ofxTxn.Amount,
			Description:    ofxTxn.Description,
			Date:           ofxTxn.Date,
			Status:         "paid",
			Tags:           pq.StringArray{},
			ExternalID:     &externalID,
			PeriodID:       &period.ID,
			ReferenceMonth: &referenceMonth,
			CreatedAt:      now,
		}

		// Inserir na transação
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
			// Verificar se é erro de violação de constraint de unicidade (duplicata)
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				// É uma duplicata esperada - fazer rollback para o savepoint e continuar
				_, rollbackErr := tx.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", savepointName))
				if rollbackErr != nil {
					tx.Rollback()
					return nil, fmt.Errorf("erro ao fazer rollback para savepoint: %w", rollbackErr)
				}
				duplicates++
				continue
			}
			// Erro dentro da transação coloca a transação em estado de falha
			// Fazer rollback e retornar erro
			tx.Rollback()
			return nil, fmt.Errorf("erro ao inserir transação: %w", err)
		}

		// Atualizar saldo da conta
		// IMPORTANTE: A matemática é a mesma para banco e cartão
		// - Expense reduz saldo (ou aumenta dívida negativa no cartão)
		// - Income aumenta saldo (ou reduz dívida negativa no cartão)
		// Para cartão de crédito, o saldo negativo representa dívida
		var balanceChange int64
		if transactionType == "expense" {
			balanceChange = -ofxTxn.Amount // Reduz saldo (ou aumenta dívida negativa)
		} else {
			balanceChange = ofxTxn.Amount // Aumenta saldo (ou reduz dívida negativa)
		}

		queryUpdateBalance := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, queryUpdateBalance, balanceChange, accountID)
		if err != nil {
			// Erro dentro da transação coloca a transação em estado de falha
			// Fazer rollback e retornar erro
			tx.Rollback()
			return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
		}

		// Liberar savepoint após sucesso
		_, err = tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName))
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("erro ao liberar savepoint: %w", err)
		}

		created++
	}

	// Commit transação
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	// Verificar se precisa atualizar saldo inicial
	if earliestDate != nil {
		// Buscar conta atualizada para verificar initial_balance_date
		updatedAccount, err := s.accountRepo.FindByID(ctx, accountID)
		if err == nil && updatedAccount != nil {
			// Se não há initial_balance_date ou a data do extrato é anterior,
			// o usuário pode querer atualizar o saldo inicial manualmente
			if updatedAccount.InitialBalanceDate == nil || earliestDate.Before(*updatedAccount.InitialBalanceDate) {
				// Potencial necessidade de atualização de saldo inicial
				// O usuário pode fazer isso manualmente através do endpoint de atualização
			}
		}
	}

	return &dto.ImportResponse{
		TotalProcessed: totalProcessed,
		Duplicates:     duplicates,
		Created:        created,
		Errors:         errorsCount,
	}, nil
}

// ImportCSV importa transações de um arquivo CSV
func (s *importService) ImportCSV(ctx context.Context, userID uuid.UUID, req *dto.ImportCSVRequest) (*dto.ImportResponse, error) {
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

	// Parse CSV
	delimiter := ","
	if req.Delimiter != "" {
		delimiter = req.Delimiter
	}

	reader := csv.NewReader(strings.NewReader(string(req.File)))
	reader.Comma = rune(delimiter[0])

	// Ler header (assumindo formato: date,amount,description)
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler header do CSV: %w", err)
	}

	// Validar formato básico
	if len(header) < 3 {
		return nil, errors.New("CSV deve ter pelo menos 3 colunas: date, amount, description")
	}

	// Iniciar transação ACID
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	var totalProcessed, duplicates, created, errorsCount int
	now := time.Now()
	var earliestDate *time.Time // Data mais antiga do extrato

	// Processar linhas
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorsCount++
			continue
		}

		if len(record) < 3 {
			errorsCount++
			continue
		}

		totalProcessed++

		// Parse da data (formato esperado: YYYY-MM-DD ou DD/MM/YYYY)
		var date time.Time
		dateStr := strings.TrimSpace(record[0])
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			// Tentar formato brasileiro
			date, err = time.Parse("02/01/2006", dateStr)
			if err != nil {
				errorsCount++
				continue
			}
		}

		// Atualizar data mais antiga
		if earliestDate == nil || date.Before(*earliestDate) {
			earliestDate = &date
		}

		// Parse do valor (assumindo formato brasileiro: R$ 1.234,56 ou 1234.56)
		amountStr := strings.TrimSpace(record[1])
		amountStr = strings.ReplaceAll(amountStr, "R$", "")
		amountStr = strings.ReplaceAll(amountStr, "$", "")
		amountStr = strings.TrimSpace(amountStr)
		amountStr = strings.ReplaceAll(amountStr, ".", "")
		amountStr = strings.ReplaceAll(amountStr, ",", ".")

		var amount float64
		_, err = fmt.Sscanf(amountStr, "%f", &amount)
		if err != nil {
			errorsCount++
			continue
		}

		// Converter para centavos
		amountInCents := int64(amount * 100)

		description := strings.TrimSpace(record[2])
		if description == "" {
			errorsCount++
			continue
		}

		// Gerar external_id para deduplicação
		externalID := s.generateExternalID(date, amountInCents, description)

		// Verificar se já existe
		existing, err := s.transactionRepo.FindByExternalID(ctx, externalID)
		if err != nil {
			errorsCount++
			continue
		}
		if existing != nil {
			duplicates++
			continue
		}

		// Determinar tipo de período baseado no tipo da conta
		periodType := "bank"
		if account.Type == "credit" {
			periodType = "credit_card"
		}

		// Usar data da transação como mês de referência (primeiro dia do mês)
		referenceMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())

		// Obter ou criar período
		period, err := s.periodService.GetOrCreatePeriod(ctx, userID, accountID, referenceMonth, periodType)
		if err != nil {
			errorsCount++
			continue
		}

		// Criar transação
		transaction := &model.Transaction{
			ID:             uuid.New(),
			UserID:         userID,
			AccountID:      accountID,
			CategoryID:     nil,
			Type:           "expense", // assumir despesa por padrão
			Amount:         amountInCents,
			Description:    description,
			Date:           date,
			Status:         "paid",
			Tags:           pq.StringArray{},
			ExternalID:     &externalID,
			PeriodID:       &period.ID,
			ReferenceMonth: &referenceMonth,
			CreatedAt:      now,
		}

		// Inserir na transação
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
			// Verificar se é erro de violação de constraint de unicidade (duplicata)
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				// É uma duplicata esperada - tratar como duplicata e continuar
				duplicates++
				continue
			}
			// Erro dentro da transação coloca a transação em estado de falha
			// Fazer rollback e retornar erro
			tx.Rollback()
			return nil, fmt.Errorf("erro ao inserir transação: %w", err)
		}

		// Atualizar saldo da conta
		// IMPORTANTE: Mesma lógica matemática para banco e cartão
		// Para cartão de crédito, o saldo negativo representa dívida
		// CSV assume expense por padrão, então subtrai do saldo
		var balanceChange int64 = -amountInCents

		queryUpdateBalance := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, queryUpdateBalance, balanceChange, accountID)
		if err != nil {
			// Erro dentro da transação coloca a transação em estado de falha
			// Fazer rollback e retornar erro
			tx.Rollback()
			return nil, fmt.Errorf("erro ao atualizar saldo da conta: %w", err)
		}

		created++
	}

	// Commit transação
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	// Verificar se precisa atualizar saldo inicial
	// Se a data mais antiga do extrato for anterior à initial_balance_date atual,
	// ou se não há initial_balance_date definido, podemos considerar atualizar
	if earliestDate != nil {
		// Buscar conta atualizada para verificar initial_balance_date
		updatedAccount, err := s.accountRepo.FindByID(ctx, accountID)
		if err == nil && updatedAccount != nil {
			// Se não há initial_balance_date ou a data do extrato é anterior,
			// o usuário pode querer atualizar o saldo inicial manualmente
			// Por enquanto, apenas logamos essa informação
			// A reconciliação automática pode ser implementada no futuro
			if updatedAccount.InitialBalanceDate == nil || earliestDate.Before(*updatedAccount.InitialBalanceDate) {
				// Potencial necessidade de atualização de saldo inicial
				// O usuário pode fazer isso manualmente através do endpoint de atualização
			}
		}
	}

	return &dto.ImportResponse{
		TotalProcessed: totalProcessed,
		Duplicates:     duplicates,
		Created:        created,
		Errors:         errorsCount,
	}, nil
}

// PreviewOFX faz preview das transações OFX sem importar
func (s *importService) PreviewOFX(ctx context.Context, userID uuid.UUID, req *dto.ImportOFXRequest) (*dto.ImportPreviewResponse, error) {
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

	// Parse do arquivo OFX
	parser := util.NewOFXParser()
	if err := parser.Parse(req.File); err != nil {
		return nil, fmt.Errorf("erro ao parsear arquivo OFX: %w", err)
	}

	ofxTransactions := parser.GetTransactions()
	if len(ofxTransactions) == 0 {
		return nil, errors.New("nenhuma transação encontrada no arquivo OFX")
	}

	// Obter tipo do arquivo OFX para validação
	ofxFileType := parser.GetFileType()

	// Validar que o tipo da conta corresponde ao tipo do arquivo OFX
	if ofxFileType == "credit_card" && account.Type != "credit" {
		return nil, fmt.Errorf("arquivo OFX é de cartão de crédito, mas a conta selecionada não é do tipo 'credit'")
	}
	if ofxFileType == "bank" && account.Type == "credit" {
		return nil, fmt.Errorf("arquivo OFX é de extrato bancário, mas a conta selecionada é do tipo 'credit'")
	}

	var transactions []*dto.TransactionDTO
	var totalTransactions, duplicates, newTransactions int

	// Processar transações OFX
	for _, ofxTxn := range ofxTransactions {
		totalTransactions++

		// Usar FITID como external_id para deduplicação
		externalID := ofxTxn.FITID

		// Verificar se já existe
		existing, err := s.transactionRepo.FindByExternalID(ctx, externalID)
		if err == nil && existing != nil {
			duplicates++
			continue
		}

		newTransactions++

		// Determinar tipo de transação baseado no TRNTYPE do OFX
		transactionType := "expense"
		if ofxTxn.Type == "CREDIT" {
			transactionType = "income"
		}

		// Determinar mês de referência (primeiro dia do mês da data)
		referenceMonth := time.Date(ofxTxn.Date.Year(), ofxTxn.Date.Month(), 1, 0, 0, 0, 0, ofxTxn.Date.Location())
		referenceMonthStr := referenceMonth.Format("2006-01")

		// Criar DTO para preview
		transactionDTO := &dto.TransactionDTO{
			ID:             uuid.New().String(),
			UserID:         userID.String(),
			AccountID:      accountID.String(),
			Type:           transactionType,
			Amount:         ofxTxn.Amount,
			Description:    ofxTxn.Description,
			Date:           ofxTxn.Date.Format("2006-01-02"),
			ReferenceMonth: &referenceMonthStr,
			Status:         "paid",
			Tags:           []string{},
			CreatedAt:      time.Now().Format(time.RFC3339),
		}

		transactions = append(transactions, transactionDTO)
	}

	return &dto.ImportPreviewResponse{
		TotalTransactions: totalTransactions,
		Duplicates:        duplicates,
		NewTransactions:   newTransactions,
		Transactions:      transactions,
	}, nil
}

// PreviewCSV faz preview das transações CSV sem importar
func (s *importService) PreviewCSV(ctx context.Context, userID uuid.UUID, req *dto.ImportCSVRequest) (*dto.ImportPreviewResponse, error) {
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

	// Parse CSV
	delimiter := ","
	if req.Delimiter != "" {
		delimiter = req.Delimiter
	}

	reader := csv.NewReader(strings.NewReader(string(req.File)))
	reader.Comma = rune(delimiter[0])

	// Ler header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler header do CSV: %w", err)
	}

	if len(header) < 3 {
		return nil, errors.New("CSV deve ter pelo menos 3 colunas: date, amount, description")
	}

	var transactions []*dto.TransactionDTO
	var totalTransactions, duplicates, newTransactions int

	// Processar linhas
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(record) < 3 {
			continue
		}

		totalTransactions++

		// Parse da data
		var date time.Time
		dateStr := strings.TrimSpace(record[0])
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			date, err = time.Parse("02/01/2006", dateStr)
			if err != nil {
				continue
			}
		}

		// Parse do valor
		amountStr := strings.TrimSpace(record[1])
		amountStr = strings.ReplaceAll(amountStr, "R$", "")
		amountStr = strings.ReplaceAll(amountStr, "$", "")
		amountStr = strings.TrimSpace(amountStr)
		amountStr = strings.ReplaceAll(amountStr, ".", "")
		amountStr = strings.ReplaceAll(amountStr, ",", ".")

		var amount float64
		_, err = fmt.Sscanf(amountStr, "%f", &amount)
		if err != nil {
			continue
		}

		amountInCents := int64(amount * 100)
		description := strings.TrimSpace(record[2])

		// Gerar external_id
		externalID := s.generateExternalID(date, amountInCents, description)

		// Verificar se já existe
		existing, err := s.transactionRepo.FindByExternalID(ctx, externalID)
		if err == nil && existing != nil {
			duplicates++
			continue
		}

		newTransactions++

		// Determinar mês de referência (primeiro dia do mês da data)
		referenceMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		referenceMonthStr := referenceMonth.Format("2006-01")

		// Criar DTO para preview
		transactionDTO := &dto.TransactionDTO{
			ID:             uuid.New().String(),
			UserID:         userID.String(),
			AccountID:      accountID.String(),
			Type:           "expense",
			Amount:         amountInCents,
			Description:    description,
			Date:           date.Format("2006-01-02"),
			ReferenceMonth: &referenceMonthStr,
			Status:         "paid",
			Tags:           []string{},
			CreatedAt:      time.Now().Format(time.RFC3339),
		}

		transactions = append(transactions, transactionDTO)
	}

	return &dto.ImportPreviewResponse{
		TotalTransactions: totalTransactions,
		Duplicates:        duplicates,
		NewTransactions:   newTransactions,
		Transactions:      transactions,
	}, nil
}
