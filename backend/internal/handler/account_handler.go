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

type AccountHandler struct {
	accountService service.AccountService
	validator      *validator.Validate
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		validator:      validator.New(),
	}
}

// CreateAccount cria uma nova conta
// @Summary Criar conta
// @Description Cria uma nova conta financeira
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAccountRequest true "Dados da conta"
// @Success 201 {object} dto.AccountDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts [post]
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	account, err := h.accountService.CreateAccount(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar conta", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusCreated, account)
}

// GetAccount obtém uma conta por ID
// @Summary Obter conta
// @Description Retorna uma conta específica
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Success 200 {object} dto.AccountDTO
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts/{id} [get]
func (h *AccountHandler) GetAccount(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), userID, accountID)
	if err != nil {
		if errors.Is(err, appError.ErrAccountNotFound) {
			return response.NotFound(c, appError.ErrAccountNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar conta", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, account)
}

// GetAccounts lista todas as contas do usuário
// @Summary Listar contas
// @Description Retorna todas as contas do usuário autenticado
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AccountDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts [get]
func (h *AccountHandler) GetAccounts(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accounts, err := h.accountService.GetAccountsByUserID(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar contas", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, accounts)
}

// UpdateAccount atualiza uma conta
// @Summary Atualizar conta
// @Description Atualiza os dados de uma conta
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Param request body dto.UpdateAccountRequest true "Dados atualizados da conta"
// @Success 200 {object} dto.AccountDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateAccountRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	account, err := h.accountService.UpdateAccount(c.Request().Context(), userID, accountID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrAccountNotFound) {
			return response.NotFound(c, appError.ErrAccountNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao atualizar conta", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, account)
}

// DeleteAccount deleta uma conta
// @Summary Deletar conta
// @Description Remove uma conta do sistema
// @Tags accounts
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.accountService.DeleteAccount(c.Request().Context(), userID, accountID); err != nil {
		if errors.Is(err, appError.ErrAccountNotFound) {
			return response.NotFound(c, appError.ErrAccountNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao deletar conta", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateInitialBalance atualiza o saldo inicial de uma conta
// @Summary Atualizar saldo inicial
// @Description Atualiza o saldo inicial de uma conta com uma data de referência, criando transação de ajuste se necessário
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Param request body dto.UpdateInitialBalanceRequest true "Dados do saldo inicial"
// @Success 200 {object} dto.AccountDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts/{id}/initial-balance [put]
func (h *AccountHandler) UpdateInitialBalance(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateInitialBalanceRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	account, err := h.accountService.UpdateInitialBalance(c.Request().Context(), userID, accountID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrAccountNotFound) {
			return response.NotFound(c, appError.ErrAccountNotFound)
		}
		if errors.Is(err, appError.ErrReferenceDateFuture) {
			return response.BadRequest(c, appError.ErrReferenceDateFuture)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao atualizar saldo inicial", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, account)
}

// RecalculateBalance recalcula o saldo de uma conta baseado no saldo inicial e nas transações
// @Summary Recalcular saldo
// @Description Recalcula o saldo de uma conta baseado no saldo inicial e em todas as transações
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Success 200 {object} dto.AccountDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/accounts/{id}/recalculate-balance [post]
func (h *AccountHandler) RecalculateBalance(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	account, err := h.accountService.RecalculateBalance(c.Request().Context(), userID, accountID)
	if err != nil {
		if errors.Is(err, appError.ErrAccountNotFound) {
			return response.NotFound(c, appError.ErrAccountNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao recalcular saldo", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, account)
}

