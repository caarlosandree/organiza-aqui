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

type CreditCardHandler struct {
	creditCardService service.CreditCardService
	validator         *validator.Validate
}

func NewCreditCardHandler(creditCardService service.CreditCardService) *CreditCardHandler {
	return &CreditCardHandler{
		creditCardService: creditCardService,
		validator:         validator.New(),
	}
}

// CreateCreditCard cria um novo cartão de crédito
// @Summary Criar cartão de crédito
// @Description Cria um novo cartão de crédito
// @Tags credit-cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCreditCardRequest true "Dados do cartão"
// @Success 201 {object} dto.CreditCardDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards [post]
func (h *CreditCardHandler) CreateCreditCard(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateCreditCardRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	creditCard, err := h.creditCardService.CreateCreditCard(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar cartão", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, creditCard)
}

// GetCreditCard obtém um cartão por ID
// @Summary Obter cartão
// @Description Retorna um cartão específico
// @Tags credit-cards
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Success 200 {object} dto.CreditCardDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id} [get]
func (h *CreditCardHandler) GetCreditCard(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	creditCard, err := h.creditCardService.GetCreditCardByID(c.Request().Context(), userID, creditCardID)
	if err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) {
			return response.NotFound(c, appError.ErrCreditCardNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, creditCard)
}

// GetCreditCards lista cartões do usuário
// @Summary Listar cartões
// @Description Retorna todos os cartões do usuário
// @Tags credit-cards
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CreditCardDTO
// @Router /api/v1/credit-cards [get]
func (h *CreditCardHandler) GetCreditCards(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCards, err := h.creditCardService.GetCreditCardsByUserID(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar cartões", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, creditCards)
}

// UpdateCreditCard atualiza um cartão
// @Summary Atualizar cartão
// @Description Atualiza os dados de um cartão
// @Tags credit-cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Param request body dto.UpdateCreditCardRequest true "Dados atualizados do cartão"
// @Success 200 {object} dto.CreditCardDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id} [put]
func (h *CreditCardHandler) UpdateCreditCard(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateCreditCardRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	creditCard, err := h.creditCardService.UpdateCreditCard(c.Request().Context(), userID, creditCardID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) {
			return response.NotFound(c, appError.ErrCreditCardNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, creditCard)
}

// DeleteCreditCard deleta um cartão
// @Summary Deletar cartão
// @Description Remove um cartão do sistema
// @Tags credit-cards
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id} [delete]
func (h *CreditCardHandler) DeleteCreditCard(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.creditCardService.DeleteCreditCard(c.Request().Context(), userID, creditCardID); err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) {
			return response.NotFound(c, appError.ErrCreditCardNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetAvailableLimit calcula o limite disponível do cartão
// @Summary Obter limite disponível
// @Description Calcula o limite disponível do cartão
// @Tags credit-cards
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Success 200 {object} map[string]int64
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id}/available-limit [get]
func (h *CreditCardHandler) GetAvailableLimit(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	availableLimit, err := h.creditCardService.GetAvailableLimit(c.Request().Context(), userID, creditCardID)
	if err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) {
			return response.NotFound(c, appError.ErrCreditCardNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, map[string]int64{
			"available_limit": availableLimit,
		})
}

// GetBillProjection projeta a fatura futura
// @Summary Projetar fatura
// @Description Projeta a fatura futura baseada em transações pendentes
// @Tags credit-cards
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Success 200 {object} dto.ProjecaoFaturaResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id}/bill-projection [get]
func (h *CreditCardHandler) GetBillProjection(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	projection, err := h.creditCardService.GetBillProjection(c.Request().Context(), userID, creditCardID)
	if err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) {
			return response.NotFound(c, appError.ErrCreditCardNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, projection)
}

// GetBills lista faturas de um cartão
// @Summary Listar faturas
// @Description Retorna todas as faturas de um cartão
// @Tags credit-cards
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Success 200 {array} dto.CreditCardBillDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id}/bills [get]
func (h *CreditCardHandler) GetBills(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	bills, err := h.creditCardService.GetBills(c.Request().Context(), userID, creditCardID)
	if err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) {
			return response.NotFound(c, appError.ErrCreditCardNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, bills)
}

// CloseBill fecha uma fatura
// @Summary Fechar fatura
// @Description Fecha uma fatura (muda status de 'open' para 'closed')
// @Tags credit-cards
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Param bill_id path string true "ID da fatura"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id}/bills/{bill_id}/close [post]
func (h *CreditCardHandler) CloseBill(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	billID, err := parseUUIDParam(c, "bill_id")
	if err != nil {
		return err
	}

	if err := h.creditCardService.CloseBill(c.Request().Context(), userID, creditCardID, billID); err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) || errors.Is(err, appError.ErrBillNotFound) {
			return response.NotFound(c, err)
		}
		if errors.Is(err, appError.ErrAlreadyClosed) || errors.Is(err, appError.ErrAlreadyPaid) {
			return response.BadRequest(c, err)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// PayBill paga uma fatura
// @Summary Pagar fatura
// @Description Paga uma fatura criando uma transação de despesa na conta bancária
// @Tags credit-cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cartão"
// @Param bill_id path string true "ID da fatura"
// @Param request body dto.PayBillRequest true "Dados do pagamento"
// @Success 200 {object} dto.CreditCardBillDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/credit-cards/{id}/bills/{bill_id}/pay [post]
func (h *CreditCardHandler) PayBill(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	creditCardID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	billID, err := parseUUIDParam(c, "bill_id")
	if err != nil {
		return err
	}

	var req dto.PayBillRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	bill, err := h.creditCardService.PayBill(c.Request().Context(), userID, creditCardID, billID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrCreditCardNotFound) || errors.Is(err, appError.ErrBillNotFound) {
			return response.NotFound(c, err)
		}
		if errors.Is(err, appError.ErrNotClosed) {
			return response.BadRequest(c, appError.ErrNotClosed)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, bill)
}

