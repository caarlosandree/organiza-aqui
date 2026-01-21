package handler

import (
	"errors"
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

type NoteHandler struct {
	noteService service.NoteService
	validator    *validator.Validate
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
		validator:    validator.New(),
	}
}

// CreateNote cria uma nova anotação
// @Summary Criar anotação
// @Description Cria uma nova anotação
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateNoteRequest true "Dados da anotação"
// @Success 201 {object} dto.NoteDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/notes [post]
func (h *NoteHandler) CreateNote(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateNoteRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	note, err := h.noteService.CreateNote(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar anotação", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, note)
}

// GetNote obtém uma anotação por ID
// @Summary Obter anotação
// @Description Retorna uma anotação específica
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da anotação"
// @Success 200 {object} dto.NoteDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/notes/{id} [get]
func (h *NoteHandler) GetNote(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	noteID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	note, err := h.noteService.GetNoteByID(c.Request().Context(), userID, noteID)
	if err != nil {
		if errors.Is(err, appError.ErrNoteNotFound) {
			return response.NotFound(c, appError.ErrNoteNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, note)
}

// GetNotes lista todas as anotações do usuário
// @Summary Listar anotações
// @Description Retorna todas as anotações do usuário com filtros opcionais
// @Tags notes
// @Produce json
// @Security BearerAuth
// @Param tag query string false "Filtrar por tag"
// @Param is_pinned query bool false "Filtrar por fixado"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {array} dto.NoteDTO
// @Router /api/v1/notes [get]
func (h *NoteHandler) GetNotes(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	filters := &dto.NoteFilters{}

	if tag := c.QueryParam("tag"); tag != "" {
		filters.Tag = &tag
	}

	if pinnedStr := c.QueryParam("is_pinned"); pinnedStr != "" {
		if pinned, err := strconv.ParseBool(pinnedStr); err == nil {
			filters.IsPinned = &pinned
		}
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	notes, err := h.noteService.GetNotesByUserID(c.Request().Context(), userID, filters)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar anotações", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, notes)
}

// UpdateNote atualiza uma anotação
// @Summary Atualizar anotação
// @Description Atualiza os dados de uma anotação
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da anotação"
// @Param request body dto.UpdateNoteRequest true "Dados atualizados"
// @Success 200 {object} dto.NoteDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/notes/{id} [put]
func (h *NoteHandler) UpdateNote(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	noteID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateNoteRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	note, err := h.noteService.UpdateNote(c.Request().Context(), userID, noteID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrNoteNotFound) {
			return response.NotFound(c, appError.ErrNoteNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, note)
}

// DeleteNote deleta uma anotação
// @Summary Deletar anotação
// @Description Remove uma anotação
// @Tags notes
// @Security BearerAuth
// @Param id path string true "ID da anotação"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/notes/{id} [delete]
func (h *NoteHandler) DeleteNote(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	noteID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.noteService.DeleteNote(c.Request().Context(), userID, noteID); err != nil {
		if errors.Is(err, appError.ErrNoteNotFound) {
			return response.NotFound(c, appError.ErrNoteNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

