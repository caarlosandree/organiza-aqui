package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type InstallmentHandler struct {
	installmentService service.InstallmentService
	validator          *validator.Validate
}

func NewInstallmentHandler(installmentService service.InstallmentService) *InstallmentHandler {
	return &InstallmentHandler{
		installmentService: installmentService,
		validator:          validator.New(),
	}
}

// CreateInstallments cria uma transação parcelada
// @Summary Criar transação parcelada
// @Description Cria uma transação com múltiplas parcelas
// @Tags installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Dados da transação (com total_installments)"
// @Success 201 {object} dto.TransactionDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/installments [post]
func (h *InstallmentHandler) CreateInstallments(c echo.Context) error {
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

	if req.TotalInstallments == nil || *req.TotalInstallments <= 1 {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	transaction, err := h.installmentService.CreateInstallments(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar parcelas", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, transaction)
}

// GetInstallments retorna todas as parcelas de uma transação mãe
// @Summary Obter parcelas
// @Description Retorna todas as parcelas de uma transação mãe
// @Tags installments
// @Produce json
// @Security BearerAuth
// @Param parent_id path string true "ID da transação mãe"
// @Success 200 {array} dto.TransactionDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/installments/{parent_id} [get]
func (h *InstallmentHandler) GetInstallments(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	parentID, err := parseUUIDParam(c, "parent_id")
	if err != nil {
		return err
	}

	installments, err := h.installmentService.GetInstallments(c.Request().Context(), userID, parentID)
	if err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, installments)
}

// UpdateInstallment atualiza uma parcela com opções de escopo
// @Summary Atualizar parcela
// @Description Atualiza uma parcela com opções: esta, esta e futuras, todas
// @Tags installments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da parcela"
// @Param request body dto.UpdateInstallmentRequest true "Dados da atualização"
// @Success 200 {object} dto.TransactionDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/installments/{id} [put]
func (h *InstallmentHandler) UpdateInstallment(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	transactionID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateInstallmentRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	transaction, err := h.installmentService.UpdateInstallment(c.Request().Context(), userID, transactionID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, transaction)
}

// CancelFutureInstallments cancela todas as parcelas futuras
// @Summary Cancelar parcelas futuras
// @Description Cancela todas as parcelas futuras de uma transação parcelada
// @Tags installments
// @Security BearerAuth
// @Param id path string true "ID da parcela atual"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/installments/{id}/cancel-future [delete]
func (h *InstallmentHandler) CancelFutureInstallments(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	transactionID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.installmentService.CancelFutureInstallments(c.Request().Context(), userID, transactionID); err != nil {
		if errors.Is(err, appError.ErrTransactionNotFound) {
			return response.NotFound(c, appError.ErrTransactionNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

