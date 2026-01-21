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

type TaskStatusHandler struct {
	taskStatusService service.TaskStatusService
	validator         *validator.Validate
}

func NewTaskStatusHandler(taskStatusService service.TaskStatusService) *TaskStatusHandler {
	return &TaskStatusHandler{
		taskStatusService: taskStatusService,
		validator:         validator.New(),
	}
}

// CreateStatus cria um novo status
// @Summary Criar status de tarefa
// @Description Cria um novo status de tarefa
// @Tags task-statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTaskStatusRequest true "Dados do status"
// @Success 201 {object} dto.TaskStatusDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/task-statuses [post]
func (h *TaskStatusHandler) CreateStatus(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateTaskStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	status, err := h.taskStatusService.CreateStatus(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar status", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusCreated, status)
}

// GetStatus obtém um status por ID
// @Summary Obter status de tarefa
// @Description Retorna um status de tarefa específico
// @Tags task-statuses
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do status"
// @Success 200 {object} dto.TaskStatusDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/task-statuses/{id} [get]
func (h *TaskStatusHandler) GetStatus(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	statusID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	status, err := h.taskStatusService.GetStatusByID(c.Request().Context(), userID, statusID)
	if err != nil {
		if errors.Is(err, appError.ErrStatusNotFound) {
			return response.NotFound(c, appError.ErrStatusNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, status)
}

// GetStatuses lista todos os status do usuário
// @Summary Listar status de tarefas
// @Description Retorna todos os status de tarefas do usuário
// @Tags task-statuses
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.TaskStatusDTO
// @Router /api/v1/task-statuses [get]
func (h *TaskStatusHandler) GetStatuses(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	statuses, err := h.taskStatusService.GetStatusesByUserID(c.Request().Context(), userID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar statuses", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, statuses)
}

// UpdateStatus atualiza um status
// @Summary Atualizar status de tarefa
// @Description Atualiza os dados de um status de tarefa
// @Tags task-statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do status"
// @Param request body dto.UpdateTaskStatusRequest true "Dados atualizados"
// @Success 200 {object} dto.TaskStatusDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/task-statuses/{id} [put]
func (h *TaskStatusHandler) UpdateStatus(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	statusID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateTaskStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	status, err := h.taskStatusService.UpdateStatus(c.Request().Context(), userID, statusID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrStatusNotFound) {
			return response.NotFound(c, appError.ErrStatusNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, status)
}

// DeleteStatus deleta um status
// @Summary Deletar status de tarefa
// @Description Remove um status de tarefa
// @Tags task-statuses
// @Security BearerAuth
// @Param id path string true "ID do status"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/task-statuses/{id} [delete]
func (h *TaskStatusHandler) DeleteStatus(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	statusID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.taskStatusService.DeleteStatus(c.Request().Context(), userID, statusID); err != nil {
		if errors.Is(err, appError.ErrStatusNotFound) {
			return response.NotFound(c, appError.ErrStatusNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// ReorderStatuses reordena os status
// @Summary Reordenar status de tarefas
// @Description Reordena os status de tarefas do usuário
// @Tags task-statuses
// @Accept json
// @Security BearerAuth
// @Param request body []string true "Array de IDs de status na ordem desejada"
// @Success 204
// @Router /api/v1/task-statuses/reorder [post]
func (h *TaskStatusHandler) ReorderStatuses(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var statusIDs []string
	if err := c.Bind(&statusIDs); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "dados inválidos",
			Message: "deve ser um array de IDs",
		})
	}

	if err := h.taskStatusService.ReorderStatuses(c.Request().Context(), userID, statusIDs); err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao reordenar statuses", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}
