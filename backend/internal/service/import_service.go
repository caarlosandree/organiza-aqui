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
	db                *sqlx.DB
	transactionRepo   repository.TransactionRepository
	accountRepo       repository.AccountRepository
	creditCardRepo    repository.CreditCardRepository
	creditCardBillRepo repository.CreditCardBillRepository
	periodService     TransactionPeriodService
}

func NewImportService(
	db *sqlx.DB,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	creditCardRepo repository.CreditCardRepository,
	creditCardBillRepo repository.CreditCardBillRepository,
	periodService TransactionPeriodService,
) ImportService {
	return &importService{
		db:                db,
		transactionRepo:   transactionRepo,
		accountRepo:       accountRepo,
		creditCardRepo:    creditCardRepo,
		creditCardBillRepo: creditCardBillRepo,
		periodService:     periodService,
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

// generateOFXExternalID gera um external_id único para transações OFX
// Combina FITID com valor, data e descrição para garantir unicidade mesmo quando
// o FITID é duplicado (caso comum em transações relacionadas como IOF + compra)
func (s *importService) generateOFXExternalID(fitid string, date time.Time, amount int64, description string) string {
	// Normalizar descrição: remover espaços extras, converter para minúsculas
	normalizedDesc := strings.ToLower(strings.TrimSpace(description))
	normalizedDesc = strings.Join(strings.Fields(normalizedDesc), " ")

	// Criar string única: FITID + data + valor + descrição normalizada
	data := fmt.Sprintf("%s|%s|%d|%s", fitid, date.Format("2006-01-02"), amount, normalizedDesc)

	// Gerar hash MD5
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// generateFileHash gera um hash MD5 do arquivo para identificar se já foi importado
func (s *importService) generateFileHash(fileContent []byte) string {
	hash := md5.Sum(fileContent)
	return hex.EncodeToString(hash[:])
}

// normalizeDescription normaliza a descrição para comparação
func (s *importService) normalizeDescription(desc string) string {
	normalized := strings.ToLower(strings.TrimSpace(desc))
	return strings.Join(strings.Fields(normalized), " ")
}

// isBillPayment identifica se uma transação é um pagamento de fatura
// baseado na descrição e tipo da transação
func (s *importService) isBillPayment(description string, transactionType string) bool {
	// Para cartões de crédito, pagamentos são do tipo CREDIT (income)
	if transactionType != "CREDIT" {
		return false
	}

	// Normalizar descrição para comparação
	normalizedDesc := s.normalizeDescription(description)

	// Padrões comuns de descrição de pagamento de fatura
	paymentPatterns := []string{
		"pagamento recebido",
		"pagamento",
		"pagto recebido",
		"pagto",
		"payment received",
		"payment",
		"pagamento da fatura",
		"pagamento fatura",
		"fatura paga",
	}

	for _, pattern := range paymentPatterns {
		if strings.Contains(normalizedDesc, pattern) {
			return true
		}
	}

	return false
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
	if ofxFileType == "credit_card" {
		// Verificar se a conta é do tipo credit OU se está vinculada a um cartão de crédito
		isCreditAccount := account.Type == "credit"
		creditCards, err := s.creditCardRepo.FindByAccountID(ctx, accountID)
		isLinkedToCreditCard := err == nil && len(creditCards) > 0

		if !isCreditAccount && !isLinkedToCreditCard {
			return nil, fmt.Errorf("arquivo OFX é de cartão de crédito, mas a conta selecionada não é do tipo 'credit' nem está vinculada a um cartão de crédito")
		}
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

	// Criar mapa de external_ids permitidos para importação (se especificado)
	allowedExternalIDs := make(map[string]bool)
	if len(req.ExternalIDs) > 0 {
		for _, id := range req.ExternalIDs {
			allowedExternalIDs[id] = true
		}
	}

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

		// Gerar external_id combinando FITID com outros campos para garantir unicidade
		// Isso resolve casos onde o FITID é duplicado (ex: IOF + compra internacional)
		externalID := s.generateOFXExternalID(ofxTxn.FITID, ofxTxn.Date, ofxTxn.Amount, ofxTxn.Description)

		// Se foi especificada uma lista de external_ids, verificar se esta transação está na lista
		if len(allowedExternalIDs) > 0 && !allowedExternalIDs[externalID] {
			// Esta transação não está na lista de permitidas, pular
			_, err = tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName))
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("erro ao liberar savepoint: %w", err)
			}
			continue
		}

		// Verificar se já existe (sempre verificar, mesmo que external_ids sejam especificados)
		// Isso garante que não importamos duplicatas mesmo se o frontend enviar IDs incorretos
		existing, err := s.transactionRepo.FindByExternalID(ctx, externalID)
		if err != nil {
			// Erro ao verificar - pular esta transação para evitar duplicatas
			errorsCount++
			_, releaseErr := tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName))
			if releaseErr != nil {
				tx.Rollback()
				return nil, fmt.Errorf("erro ao liberar savepoint: %w", releaseErr)
			}
			continue
		}
		if existing != nil {
			// Transação já existe - marcar como duplicata e pular
			duplicates++
			_, releaseErr := tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName))
			if releaseErr != nil {
				tx.Rollback()
				return nil, fmt.Errorf("erro ao liberar savepoint: %w", releaseErr)
			}
			continue
		}

		// Determinar tipo de transação baseado no TRNTYPE do OFX
		transactionType := "expense"
		if ofxTxn.Type == "CREDIT" {
			transactionType = "income"
		}

		// Verificar se é um pagamento de fatura
		isPayment := s.isBillPayment(ofxTxn.Description, ofxTxn.Type)

		// Determinar mês de referência e fatura relacionada
		var referenceMonth time.Time
		var billID *uuid.UUID

		if isPayment && ofxFileType == "credit_card" {
			// Para pagamentos de fatura, determinar a fatura correta
			creditCards, err := s.creditCardRepo.FindByAccountID(ctx, accountID)
			if err == nil && len(creditCards) > 0 {
				// Usar o primeiro cartão vinculado (assumindo um cartão por conta)
				creditCard := creditCards[0]

				// Determinar qual fatura o pagamento pertence
				// Se o pagamento é antes do dia de fechamento, pertence ao mês anterior
				paymentDate := ofxTxn.Date
				paymentYear := paymentDate.Year()
				paymentMonth := int(paymentDate.Month())
				paymentDay := paymentDate.Day()

				// Se o pagamento é antes do dia de fechamento, pertence ao mês anterior
				if paymentDay < creditCard.ClosingDay {
					// Pagamento pertence ao mês anterior
					paymentMonth--
					if paymentMonth < 1 {
						paymentMonth = 12
						paymentYear--
					}
				}

				// Buscar fatura do mês determinado diretamente na transação
				// IMPORTANTE: Não criar a fatura se não existir
				// Apenas relacionar se já existir
				var bill model.CreditCardBill
				billQuery := `SELECT id, credit_card_id, month, year, status, closing_date, due_date, payment_transaction_id, created_at, updated_at
				              FROM credit_card_bills 
				              WHERE credit_card_id = $1 AND year = $2 AND month = $3`
				err = tx.GetContext(ctx, &bill, billQuery, creditCard.ID, paymentYear, paymentMonth)
				if err == nil {
					// Fatura existe - relacionar o pagamento
					billID = &bill.ID
				}
				// Se a fatura não existir (sql.ErrNoRows), billID permanece nil
				// A transação será criada normalmente, mas sem relacionamento com fatura
				// Quando o usuário importar o extrato do mês anterior, a fatura será criada
				// e o pagamento poderá ser relacionado depois

				// Mesmo que seja pagamento, usar o mês da fatura como referência
				referenceMonth = time.Date(paymentYear, time.Month(paymentMonth), 1, 0, 0, 0, 0, ofxTxn.Date.Location())
			} else {
				// Não encontrou cartão - usar data da transação como referência
				referenceMonth = time.Date(ofxTxn.Date.Year(), ofxTxn.Date.Month(), 1, 0, 0, 0, 0, ofxTxn.Date.Location())
			}
		} else {
			// Para transações normais, usar data da transação como mês de referência
			referenceMonth = time.Date(ofxTxn.Date.Year(), ofxTxn.Date.Month(), 1, 0, 0, 0, 0, ofxTxn.Date.Location())
		}

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

		// Se é um pagamento e encontramos uma fatura existente, relacionar
		if isPayment && billID != nil {
			// Atualizar fatura com payment_transaction_id
			// Apenas atualizar se ainda não tiver payment_transaction_id (evitar sobrescrever)
			updateBillQuery := `
				UPDATE credit_card_bills 
				SET payment_transaction_id = $1, status = 'paid', updated_at = $2
				WHERE id = $3 AND payment_transaction_id IS NULL
			`
			_, err = tx.ExecContext(ctx, updateBillQuery, transaction.ID, now, *billID)
			if err != nil {
				// Erro ao atualizar fatura - logar mas não falhar a importação
				// A transação já foi criada com sucesso
			}
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
	if ofxFileType == "credit_card" {
		// Verificar se a conta é do tipo credit OU se está vinculada a um cartão de crédito
		isCreditAccount := account.Type == "credit"
		creditCards, err := s.creditCardRepo.FindByAccountID(ctx, accountID)
		isLinkedToCreditCard := err == nil && len(creditCards) > 0

		if !isCreditAccount && !isLinkedToCreditCard {
			return nil, fmt.Errorf("arquivo OFX é de cartão de crédito, mas a conta selecionada não é do tipo 'credit' nem está vinculada a um cartão de crédito")
		}
	}
	if ofxFileType == "bank" && account.Type == "credit" {
		return nil, fmt.Errorf("arquivo OFX é de extrato bancário, mas a conta selecionada é do tipo 'credit'")
	}

	// Gerar hash do arquivo para identificar se já foi importado
	fileHash := s.generateFileHash(req.File)

	var transactions []*dto.TransactionPreviewDTO
	var totalTransactions, duplicates, newTransactions int

	// Processar transações OFX
	for _, ofxTxn := range ofxTransactions {
		totalTransactions++

		// Gerar external_id combinando FITID com outros campos para garantir unicidade
		// Isso resolve casos onde o FITID é duplicado (ex: IOF + compra internacional)
		externalID := s.generateOFXExternalID(ofxTxn.FITID, ofxTxn.Date, ofxTxn.Amount, ofxTxn.Description)

		// Verificar se já existe
		existing, err := s.transactionRepo.FindByExternalID(ctx, externalID)
		status := "new"
		if err == nil && existing != nil {
			duplicates++
			status = "existing"
		} else {
			newTransactions++
		}

		// Determinar tipo de transação baseado no TRNTYPE do OFX
		transactionType := "expense"
		if ofxTxn.Type == "CREDIT" {
			transactionType = "income"
		}

		// Determinar mês de referência (primeiro dia do mês da data)
		referenceMonth := time.Date(ofxTxn.Date.Year(), ofxTxn.Date.Month(), 1, 0, 0, 0, 0, ofxTxn.Date.Location())
		referenceMonthStr := referenceMonth.Format("2006-01")

		// Criar DTO para preview com status
		transactionPreviewDTO := &dto.TransactionPreviewDTO{
			TransactionDTO: dto.TransactionDTO{
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
			},
			ExternalID: externalID,
			Status:     status,
		}

		transactions = append(transactions, transactionPreviewDTO)
	}

	return &dto.ImportPreviewResponse{
		FileHash:          fileHash,
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

	// Gerar hash do arquivo para identificar se já foi importado
	fileHash := s.generateFileHash(req.File)

	var transactions []*dto.TransactionPreviewDTO
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
		status := "new"
		if err == nil && existing != nil {
			duplicates++
			status = "existing"
		} else {
			newTransactions++
		}

		// Determinar mês de referência (primeiro dia do mês da data)
		referenceMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		referenceMonthStr := referenceMonth.Format("2006-01")

		// Criar DTO para preview com status
		transactionPreviewDTO := &dto.TransactionPreviewDTO{
			TransactionDTO: dto.TransactionDTO{
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
			},
			ExternalID: externalID,
			Status:     status,
		}

		transactions = append(transactions, transactionPreviewDTO)
	}

	return &dto.ImportPreviewResponse{
		FileHash:          fileHash,
		TotalTransactions: totalTransactions,
		Duplicates:        duplicates,
		NewTransactions:   newTransactions,
		Transactions:      transactions,
	}, nil
}
