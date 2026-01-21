package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetIncomeExpenseByPeriod obtém receitas e despesas por período
// @Summary Receitas e Despesas por Período
// @Description Retorna o total de receitas e despesas em um período
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Success 200 {object} dto.IncomeExpenseByPeriodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/analytics/income-expense [get]
func (h *AnalyticsHandler) GetIncomeExpenseByPeriod(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
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

	result, err := h.analyticsService.GetIncomeExpenseByPeriod(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar analytics", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// GetCategoryBreakdown obtém breakdown por categoria
// @Summary Breakdown por Categoria
// @Description Retorna o total gasto/recebido por categoria
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Param type query string true "Tipo (income ou expense)"
// @Success 200 {array} dto.CategoryBreakdownDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/analytics/category-breakdown [get]
func (h *AnalyticsHandler) GetCategoryBreakdown(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	transactionType := c.QueryParam("type")

	if startDateStr == "" || endDateStr == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if transactionType != "income" && transactionType != "expense" {
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

	result, err := h.analyticsService.GetCategoryBreakdown(c.Request().Context(), userID, startDate, endDate, transactionType)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar breakdown", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// GetMonthlyTrend obtém tendência mensal
// @Summary Tendência Mensal
// @Description Retorna receitas e despesas por mês
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Param months query int false "Número de meses (padrão: 6)"
// @Success 200 {array} dto.MonthlyTrendDTO
// @Router /api/v1/analytics/monthly-trend [get]
func (h *AnalyticsHandler) GetMonthlyTrend(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	monthsStr := c.QueryParam("months")
	months := 6 // padrão
	if monthsStr != "" {
		if parsed, err := strconv.Atoi(monthsStr); err == nil && parsed > 0 {
			months = parsed
		}
	}

	result, err := h.analyticsService.GetMonthlyTrend(c.Request().Context(), userID, months)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar tendência mensal", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// GetPatrimonioLiquido calcula o patrimônio líquido
// @Summary Patrimônio Líquido
// @Description Calcula o patrimônio líquido (soma de saldos de todas contas - faturas abertas de cartões)
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PatrimonioLiquidoResponse
// @Router /api/v1/analytics/patrimonio-liquido [get]
func (h *AnalyticsHandler) GetPatrimonioLiquido(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	result, err := h.analyticsService.GetPatrimonioLiquido(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao calcular patrimônio líquido", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// GetCalendarioVencimentos retorna o calendário de vencimentos
// @Summary Calendário de Vencimentos
// @Description Retorna transações e faturas com status 'pending' ordenadas por data
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Param days query int false "Número de dias à frente (padrão: 30)"
// @Success 200 {object} dto.CalendarioVencimentosResponse
// @Router /api/v1/analytics/calendario-vencimentos [get]
func (h *AnalyticsHandler) GetCalendarioVencimentos(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	daysStr := c.QueryParam("days")
	days := 30 // padrão
	if daysStr != "" {
		if parsed, err := strconv.Atoi(daysStr); err == nil && parsed > 0 {
			days = parsed
		}
	}

	result, err := h.analyticsService.GetCalendarioVencimentos(c.Request().Context(), userID, days)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar calendário de vencimentos", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// GetGastosPorTag agrupa gastos por tag
// @Summary Gastos por Tag
// @Description Retorna gastos agrupados por tag
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Data inicial (YYYY-MM-DD)"
// @Param end_date query string true "Data final (YYYY-MM-DD)"
// @Success 200 {object} dto.GastosPorTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/analytics/gastos-por-tag [get]
func (h *AnalyticsHandler) GetGastosPorTag(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
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

	result, err := h.analyticsService.GetGastosPorTag(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar gastos por tag", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

