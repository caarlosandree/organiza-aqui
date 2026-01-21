package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type TransactionPeriodHandler struct {
	periodService service.TransactionPeriodService
	validator     *validator.Validate
}

func NewTransactionPeriodHandler(periodService service.TransactionPeriodService) *TransactionPeriodHandler {
	return &TransactionPeriodHandler{
		periodService: periodService,
		validator:     validator.New(),
	}
}

// GetTransactionPeriods lista períodos do usuário
// @Summary Listar períodos de transações
// @Description Lista todos os períodos de transações do usuário com filtros opcionais
// @Tags transaction-periods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account_id query string false "ID da conta"
// @Param period_type query string false "Tipo do período (bank ou credit_card)"
// @Param year query int false "Ano"
// @Param month query int false "Mês (1-12)"
// @Param status query string false "Status (open, closed, archived)"
// @Param reference_month query string false "Mês de referência (YYYY-MM)"
// @Success 200 {array} dto.TransactionPeriodDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/transaction-periods [get]
func (h *TransactionPeriodHandler) GetTransactionPeriods(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	filters := &dto.TransactionPeriodFilters{}

	if accountID := c.QueryParam("account_id"); accountID != "" {
		filters.AccountID = &accountID
	}

	if periodType := c.QueryParam("period_type"); periodType != "" {
		if periodType != "bank" && periodType != "credit_card" {
		return response.BadRequest(c, appError.ErrInvalidInput)
		}
		filters.PeriodType = &periodType
	}

	if yearStr := c.QueryParam("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
		}
		filters.Year = &year
	}

	if monthStr := c.QueryParam("month"); monthStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
		return response.BadRequest(c, appError.ErrInvalidInput)
		}
		filters.Month = &month
	}

	if status := c.QueryParam("status"); status != "" {
		if status != "open" && status != "closed" && status != "archived" {
		return response.BadRequest(c, appError.ErrInvalidInput)
		}
		filters.Status = &status
	}

	if referenceMonth := c.QueryParam("reference_month"); referenceMonth != "" {
		filters.ReferenceMonth = &referenceMonth
	}

	periods, err := h.periodService.GetPeriodsByUserID(c.Request().Context(), userID, filters)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar períodos", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, periods)
}

// GetTransactionPeriod obtém um período específico com suas transações
// @Summary Obter período com transações
// @Description Retorna um período específico com todas as suas transações e estatísticas
// @Tags transaction-periods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do período"
// @Success 200 {object} dto.PeriodWithTransactionsDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/transaction-periods/:id [get]
func (h *TransactionPeriodHandler) GetTransactionPeriod(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	periodID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	period, err := h.periodService.GetPeriodWithTransactions(c.Request().Context(), userID, periodID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar período", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, period)
}

