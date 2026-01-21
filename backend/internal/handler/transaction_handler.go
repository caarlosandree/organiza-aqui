package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	validator          *validator.Validate
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		validator:          validator.New(),
	}
}

// CreateTransaction cria uma nova transação
// @Summary Criar transação
// @Description Cria uma nova transação financeira
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Dados da transação"
// @Success 201 {object} dto.TransactionDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	transaction, err := h.transactionService.CreateTransaction(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar transação", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, transaction)
}

// GetTransaction obtém uma transação por ID
// @Summary Obter transação
// @Description Retorna uma transação específica
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da transação"
// @Success 200 {object} dto.TransactionDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	transaction, err := h.transactionService.GetTransactionByID(c.Request().Context(), userID, transactionID)
	if err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, transaction)
}

// GetTransactions lista transações do usuário
// @Summary Listar transações
// @Description Retorna transações do usuário autenticado com filtros opcionais
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param account_id query string false "Filtrar por conta"
// @Param category_id query string false "Filtrar por categoria"
// @Param type query string false "Filtrar por tipo (income, expense, transfer)"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	filters := &dto.TransactionFilters{}

	if accountID := c.QueryParam("account_id"); accountID != "" {
		filters.AccountID = &accountID
	}

	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		filters.CategoryID = &categoryID
	}

	if transactionType := c.QueryParam("type"); transactionType != "" {
		filters.Type = &transactionType
	}

	if status := c.QueryParam("status"); status != "" {
		filters.Status = &status
	}

	if parentTransactionID := c.QueryParam("parent_transaction_id"); parentTransactionID != "" {
		filters.ParentTransactionID = &parentTransactionID
	}

	if startDate := c.QueryParam("start_date"); startDate != "" {
		filters.StartDate = &startDate
	}

	if endDate := c.QueryParam("end_date"); endDate != "" {
		filters.EndDate = &endDate
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	transactions, total, err := h.transactionService.GetTransactionsByUserID(c.Request().Context(), userID, filters)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar transações", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, map[string]interface{}{
		"data":  transactions,
		"total": total,
	})
}

// UpdateTransaction atualiza uma transação
// @Summary Atualizar transação
// @Description Atualiza os dados de uma transação
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da transação"
// @Param request body dto.UpdateTransactionRequest true "Dados atualizados da transação"
// @Success 200 {object} dto.TransactionDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	var req dto.UpdateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	transaction, err := h.transactionService.UpdateTransaction(c.Request().Context(), userID, transactionID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, transaction)
}

// DeleteTransaction deleta uma transação
// @Summary Deletar transação
// @Description Remove uma transação do sistema
// @Tags transactions
// @Security BearerAuth
// @Param id path string true "ID da transação"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.transactionService.DeleteTransaction(c.Request().Context(), userID, transactionID); err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetStatement obtém o extrato de uma conta
// @Summary Obter extrato
// @Description Retorna o extrato de uma conta para um período específico
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param account_id query string true "ID da conta"
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Success 200 {object} dto.StatementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/transactions/statement [get]
func (h *TransactionHandler) GetStatement(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountIDStr := c.QueryParam("account_id")
	if accountIDStr == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	if startDateStr == "" || endDateStr == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if startDate.After(endDate) {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	statement, err := h.transactionService.GetStatement(c.Request().Context(), userID, accountID, startDate, endDate)
	if err != nil {
		if errors.Is(err, appError.ErrAccountNotFound) {
			return response.NotFound(c, appError.ErrAccountNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao gerar extrato", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, statement)
}

// UpdateTransactionStatus atualiza o status de uma transação
// @Summary Atualizar status da transação
// @Description Atualiza o status de uma transação (pending, paid, cancelled)
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da transação"
// @Param status body map[string]string true "Status da transação" SchemaExample({"status": "paid"})
// @Success 200 {object} dto.TransactionDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/transactions/{id}/status [patch]
func (h *TransactionHandler) UpdateTransactionStatus(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=pending paid cancelled"`
	}
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	// Buscar transação atual
	transaction, err := h.transactionService.GetTransactionByID(c.Request().Context(), userID, transactionID)
	if err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	// Atualizar apenas o status
	updateReq := dto.UpdateTransactionRequest{
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date,
		Status:      req.Status,
		Tags:        transaction.Tags,
	}

	updated, err := h.transactionService.UpdateTransaction(c.Request().Context(), userID, transactionID, &updateReq)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, updated)
}
