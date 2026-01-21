package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type BankHandler struct {
	bankService service.BankService
	validator   *validator.Validate
}

func NewBankHandler(bankService service.BankService) *BankHandler {
	return &BankHandler{
		bankService: bankService,
		validator:   validator.New(),
	}
}

// GetAllBanks obtém todos os bancos
// @Summary Lista todos os bancos
// @Description Retorna uma lista com todos os bancos cadastrados
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.BankDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/banks [get]
func (h *BankHandler) GetAllBanks(c echo.Context) error {
	banks, err := h.bankService.GetAllBanks(c.Request().Context())
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar bancos", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, banks)
}

// GetBank obtém um banco por ID
// @Summary Obter banco
// @Description Retorna um banco específico
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do banco"
// @Success 200 {object} dto.BankDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/banks/{id} [get]
func (h *BankHandler) GetBank(c echo.Context) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	bank, err := h.bankService.GetBankByID(c.Request().Context(), id)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar banco", err)
		if errors.Is(err, appError.ErrNotFound) {
			return response.NotFound(c, appError.ErrNotFound)
		}
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, bank)
}

// SyncBanks sincroniza bancos com a API BrasilAPI
// @Summary Sincronizar bancos
// @Description Sincroniza a lista de bancos com a API BrasilAPI
// @Tags banks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/banks/sync [post]
func (h *BankHandler) SyncBanks(c echo.Context) error {
	if err := h.bankService.SyncBanks(c.Request().Context()); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao sincronizar bancos", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, map[string]string{
			"message": "Sincronização concluída com sucesso",
		})
}
