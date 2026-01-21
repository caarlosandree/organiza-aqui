package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	"github.com/organiza-aqui/backend/pkg/response"
)

type ImportHandler struct {
	importService service.ImportService
	validator     *validator.Validate
}

func NewImportHandler(importService service.ImportService) *ImportHandler {
	return &ImportHandler{
		importService: importService,
		validator:     validator.New(),
	}
}

// ImportOFX importa transações de um arquivo OFX
// @Summary Importar OFX
// @Description Importa transações de um arquivo OFX. Se external_ids for fornecido, importa apenas as transações especificadas.
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param account_id formData string true "ID da conta"
// @Param file formData file true "Arquivo OFX"
// @Param external_ids formData string false "Lista de external_ids para importar apenas transações específicas (JSON array)"
// @Success 200 {object} dto.ImportResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/import/ofx [post]
func (h *ImportHandler) ImportOFX(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID := c.FormValue("account_id")
	if accountID == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	src, err := file.Open()
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	if _, err := src.Read(fileBytes); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	req := &dto.ImportOFXRequest{
		AccountID: accountID,
		File:       fileBytes,
	}

	// Tentar obter external_ids do form-data (opcional)
	externalIDsStr := c.FormValue("external_ids")
	if externalIDsStr != "" {
		var externalIDs []string
		if err := json.Unmarshal([]byte(externalIDsStr), &externalIDs); err != nil {
			// Se não for JSON válido, tentar como string única
			externalIDs = []string{externalIDsStr}
		}
		req.ExternalIDs = externalIDs
	}

	result, err := h.importService.ImportOFX(c.Request().Context(), userID, req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao importar OFX", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// ImportCSV importa transações de um arquivo CSV
// @Summary Importar CSV
// @Description Importa transações de um arquivo CSV
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param account_id formData string true "ID da conta"
// @Param file formData file true "Arquivo CSV"
// @Param delimiter formData string false "Delimitador (padrão: ,)"
// @Success 200 {object} dto.ImportResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/import/csv [post]
func (h *ImportHandler) ImportCSV(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID := c.FormValue("account_id")
	if accountID == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	src, err := file.Open()
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	if _, err := src.Read(fileBytes); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	req := &dto.ImportCSVRequest{
		AccountID: accountID,
		File:       fileBytes,
		Delimiter:  c.FormValue("delimiter"),
	}

	result, err := h.importService.ImportCSV(c.Request().Context(), userID, req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao importar CSV", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// PreviewOFX faz preview das transações OFX sem importar
// @Summary Preview OFX
// @Description Faz preview das transações de um arquivo OFX sem importar
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param account_id formData string true "ID da conta"
// @Param file formData file true "Arquivo OFX"
// @Success 200 {object} dto.ImportPreviewResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/import/ofx/preview [post]
func (h *ImportHandler) PreviewOFX(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID := c.FormValue("account_id")
	if accountID == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	src, err := file.Open()
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	if _, err := src.Read(fileBytes); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	req := &dto.ImportOFXRequest{
		AccountID: accountID,
		File:       fileBytes,
	}

	result, err := h.importService.PreviewOFX(c.Request().Context(), userID, req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

// PreviewCSV faz preview das transações CSV sem importar
// @Summary Preview CSV
// @Description Faz preview das transações de um arquivo CSV sem importar
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param account_id formData string true "ID da conta"
// @Param file formData file true "Arquivo CSV"
// @Param delimiter formData string false "Delimitador (padrão: ,)"
// @Success 200 {object} dto.ImportPreviewResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/import/csv/preview [post]
func (h *ImportHandler) PreviewCSV(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	accountID := c.FormValue("account_id")
	if accountID == "" {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	src, err := file.Open()
	if err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	if _, err := src.Read(fileBytes); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	req := &dto.ImportCSVRequest{
		AccountID: accountID,
		File:       fileBytes,
		Delimiter:  c.FormValue("delimiter"),
	}

	result, err := h.importService.PreviewCSV(c.Request().Context(), userID, req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, result)
}

