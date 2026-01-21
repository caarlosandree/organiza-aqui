package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/organiza-aqui/backend/internal/dto"
	"github.com/organiza-aqui/backend/pkg/logger"
	"go.uber.org/zap"
)

const (
	BrasilAPIBaseURL = "https://brasilapi.com.br/api/banks/v1"
	RequestTimeout   = 30 * time.Second
)

// BrasilAPIClient é o cliente HTTP para a API BrasilAPI
type BrasilAPIClient struct {
	httpClient *http.Client
}

// NewBrasilAPIClient cria uma nova instância do cliente BrasilAPI
func NewBrasilAPIClient() *BrasilAPIClient {
	return &BrasilAPIClient{
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

// GetBanks busca todos os bancos da API BrasilAPI
func (c *BrasilAPIClient) GetBanks() ([]dto.BrasilAPIBankResponse, error) {
	req, err := http.NewRequest("GET", BrasilAPIBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "organiza-aqui/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Erro ao fazer requisição para BrasilAPI", err)
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("erro na resposta da API: status %d", resp.StatusCode)
		logger.Error("Erro na resposta da BrasilAPI", err,
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var banks []dto.BrasilAPIBankResponse
	if err := json.Unmarshal(body, &banks); err != nil {
		logger.Error("Erro ao fazer parse da resposta", err,
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	return banks, nil
}
