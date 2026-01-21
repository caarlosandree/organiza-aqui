package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// OFXTransaction representa uma transação extraída do arquivo OFX
type OFXTransaction struct {
	Type        string    // "DEBIT" ou "CREDIT"
	Date        time.Time // Data da transação
	Amount      int64     // Valor em centavos (sempre positivo)
	Description string    // Descrição/Memo
	FITID       string    // Financial Institution Transaction ID (para deduplicação)
}

// OFXParser é responsável por fazer parse de arquivos OFX
type OFXParser struct {
	transactions []OFXTransaction
}

// NewOFXParser cria um novo parser OFX
func NewOFXParser() *OFXParser {
	return &OFXParser{
		transactions: make([]OFXTransaction, 0),
	}
}

// Parse faz o parse do conteúdo do arquivo OFX
func (p *OFXParser) Parse(content []byte) error {
	contentStr := string(content)

	// Verificar se é formato OFX válido
	if !strings.Contains(contentStr, "<OFX>") {
		return fmt.Errorf("arquivo OFX inválido: não contém tag <OFX>")
	}

	// Extrair seção de transações de cartão de crédito
	if strings.Contains(contentStr, "<CREDITCARDMSGSRSV1>") {
		if err := p.parseCreditCardTransactions(contentStr); err != nil {
			return fmt.Errorf("erro ao parsear transações de cartão de crédito: %w", err)
		}
	}

	// Extrair seção de transações bancárias
	if strings.Contains(contentStr, "<BANKMSGSRSV1>") {
		if err := p.parseBankTransactions(contentStr); err != nil {
			return fmt.Errorf("erro ao parsear transações bancárias: %w", err)
		}
	}

	if len(p.transactions) == 0 {
		return fmt.Errorf("nenhuma transação encontrada no arquivo OFX")
	}

	return nil
}

// GetTransactions retorna as transações parseadas
func (p *OFXParser) GetTransactions() []OFXTransaction {
	return p.transactions
}

// parseCreditCardTransactions faz parse das transações de cartão de crédito
func (p *OFXParser) parseCreditCardTransactions(content string) error {
	// Encontrar a seção CREDITCARDMSGSRSV1
	creditCardStart := strings.Index(content, "<CREDITCARDMSGSRSV1>")
	if creditCardStart == -1 {
		return fmt.Errorf("seção CREDITCARDMSGSRSV1 não encontrada")
	}

	// Encontrar o final da seção CREDITCARDMSGSRSV1
	creditCardEnd := strings.Index(content[creditCardStart:], "</CREDITCARDMSGSRSV1>")
	if creditCardEnd == -1 {
		return fmt.Errorf("tag de fechamento CREDITCARDMSGSRSV1 não encontrada")
	}
	creditCardSection := content[creditCardStart : creditCardStart+creditCardEnd+len("</CREDITCARDMSGSRSV1>")]

	// Encontrar a seção BANKTRANLIST dentro de CREDITCARDMSGSRSV1
	startTag := "<BANKTRANLIST>"
	endTag := "</BANKTRANLIST>"

	startIdx := strings.Index(creditCardSection, startTag)
	if startIdx == -1 {
		return fmt.Errorf("seção BANKTRANLIST não encontrada em CREDITCARDMSGSRSV1")
	}

	endIdx := strings.Index(creditCardSection[startIdx:], endTag)
	if endIdx == -1 {
		return fmt.Errorf("tag de fechamento BANKTRANLIST não encontrada")
	}

	transactionList := creditCardSection[startIdx : startIdx+endIdx+len(endTag)]

	// Extrair todas as transações STMTTRN
	transactions := p.extractSTMTTRN(transactionList)

	for _, txn := range transactions {
		parsed, err := p.parseSTMTTRN(txn)
		if err != nil {
			// Log erro mas continue processando outras transações
			continue
		}
		p.transactions = append(p.transactions, *parsed)
	}

	return nil
}

// parseBankTransactions faz parse das transações bancárias
func (p *OFXParser) parseBankTransactions(content string) error {
	// Encontrar a seção BANKMSGSRSV1
	bankStart := strings.Index(content, "<BANKMSGSRSV1>")
	if bankStart == -1 {
		return fmt.Errorf("seção BANKMSGSRSV1 não encontrada")
	}

	// Encontrar o final da seção BANKMSGSRSV1
	bankEnd := strings.Index(content[bankStart:], "</BANKMSGSRSV1>")
	if bankEnd == -1 {
		return fmt.Errorf("tag de fechamento BANKMSGSRSV1 não encontrada")
	}
	bankSection := content[bankStart : bankStart+bankEnd+len("</BANKMSGSRSV1>")]

	// Encontrar a seção BANKTRANLIST dentro de BANKMSGSRSV1
	startTag := "<BANKTRANLIST>"
	endTag := "</BANKTRANLIST>"

	startIdx := strings.Index(bankSection, startTag)
	if startIdx == -1 {
		return fmt.Errorf("seção BANKTRANLIST não encontrada em BANKMSGSRSV1")
	}

	endIdx := strings.Index(bankSection[startIdx:], endTag)
	if endIdx == -1 {
		return fmt.Errorf("tag de fechamento BANKTRANLIST não encontrada")
	}

	transactionList := bankSection[startIdx : startIdx+endIdx+len(endTag)]

	// Extrair todas as transações STMTTRN
	transactions := p.extractSTMTTRN(transactionList)

	for _, txn := range transactions {
		parsed, err := p.parseSTMTTRN(txn)
		if err != nil {
			// Log erro mas continue processando outras transações
			continue
		}
		p.transactions = append(p.transactions, *parsed)
	}

	return nil
}

// extractSTMTTRN extrai todas as tags STMTTRN do conteúdo
func (p *OFXParser) extractSTMTTRN(content string) []string {
	var transactions []string
	startTag := "<STMTTRN>"
	endTag := "</STMTTRN>"

	idx := 0
	for {
		startIdx := strings.Index(content[idx:], startTag)
		if startIdx == -1 {
			break
		}
		startIdx += idx

		endIdx := strings.Index(content[startIdx:], endTag)
		if endIdx == -1 {
			break
		}
		endIdx += startIdx + len(endTag)

		transactions = append(transactions, content[startIdx:endIdx])
		idx = endIdx
	}

	return transactions
}

// parseSTMTTRN faz parse de uma tag STMTTRN individual
func (p *OFXParser) parseSTMTTRN(content string) (*OFXTransaction, error) {
	txn := &OFXTransaction{}

	// Extrair TRNTYPE
	if trnType := p.extractTag(content, "TRNTYPE"); trnType != "" {
		txn.Type = strings.TrimSpace(trnType)
	} else {
		return nil, fmt.Errorf("TRNTYPE não encontrado")
	}

	// Extrair DTPOSTED
	dtPosted := p.extractTag(content, "DTPOSTED")
	if dtPosted == "" {
		return nil, fmt.Errorf("DTPOSTED não encontrado")
	}

	date, err := p.parseOFXDate(dtPosted)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear data: %w", err)
	}
	txn.Date = date

	// Extrair TRNAMT
	trnAmt := p.extractTag(content, "TRNAMT")
	if trnAmt == "" {
		return nil, fmt.Errorf("TRNAMT não encontrado")
	}

	amount, err := p.parseAmount(trnAmt)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear valor: %w", err)
	}
	txn.Amount = amount

	// Extrair MEMO (descrição)
	memo := p.extractTag(content, "MEMO")
	txn.Description = strings.TrimSpace(memo)
	if txn.Description == "" {
		txn.Description = "Transação sem descrição"
	}

	// Extrair FITID (obrigatório para deduplicação)
	fitid := p.extractTag(content, "FITID")
	if fitid == "" {
		return nil, fmt.Errorf("FITID não encontrado")
	}
	txn.FITID = strings.TrimSpace(fitid)

	return txn, nil
}

// extractTag extrai o conteúdo de uma tag XML/SGML
func (p *OFXParser) extractTag(content, tagName string) string {
	startTag := "<" + tagName + ">"
	endTag := "</" + tagName + ">"

	startIdx := strings.Index(content, startTag)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(startTag)

	endIdx := strings.Index(content[startIdx:], endTag)
	if endIdx == -1 {
		return ""
	}

	return content[startIdx : startIdx+endIdx]
}

// parseOFXDate faz parse de uma data no formato OFX
// Formato: YYYYMMDDHHMMSS[.fff][TZ]
// Exemplo: 20260118000000[-3:BRT] ou 20260120192124[0:GMT]
func (p *OFXParser) parseOFXDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Remover timezone se presente (ex: [-3:BRT] ou [0:GMT])
	re := regexp.MustCompile(`\[.*?\]`)
	dateStr = re.ReplaceAllString(dateStr, "")

	// Formato mínimo: YYYYMMDD (8 caracteres)
	if len(dateStr) < 8 {
		return time.Time{}, fmt.Errorf("data OFX muito curta: %s", dateStr)
	}

	// Extrair apenas a parte da data (YYYYMMDD)
	datePart := dateStr[:8]

	// Parse da data
	year, err := strconv.Atoi(datePart[:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("erro ao parsear ano: %w", err)
	}

	month, err := strconv.Atoi(datePart[4:6])
	if err != nil {
		return time.Time{}, fmt.Errorf("erro ao parsear mês: %w", err)
	}

	day, err := strconv.Atoi(datePart[6:8])
	if err != nil {
		return time.Time{}, fmt.Errorf("erro ao parsear dia: %w", err)
	}

	// Criar time.Time (usar UTC como padrão, já que o OFX pode não especificar timezone)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// parseAmount faz parse de um valor monetário e converte para centavos
// Valores negativos indicam débito, positivos indicam crédito
// Retorna sempre valor positivo em centavos
func (p *OFXParser) parseAmount(amountStr string) (int64, error) {
	amountStr = strings.TrimSpace(amountStr)

	// Remover espaços e caracteres não numéricos (exceto ponto e sinal de menos)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")

	// Parse do valor
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao parsear valor: %w", err)
	}

	// Converter para centavos (sempre positivo)
	amountInCents := int64(amount * 100)
	if amountInCents < 0 {
		amountInCents = -amountInCents
	}

	return amountInCents, nil
}
