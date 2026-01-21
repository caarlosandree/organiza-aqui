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

type TaskHandler struct {
	taskService service.TaskService
	validator   *validator.Validate
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		validator:   validator.New(),
	}
}

// CreateTask cria uma nova tarefa
// @Summary Criar tarefa
// @Description Cria uma nova tarefa
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTaskRequest true "Dados da tarefa"
// @Success 201 {object} dto.TaskDTO
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	task, err := h.taskService.CreateTask(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusCreated, task)
}

// GetTask obtém uma tarefa por ID
// @Summary Obter tarefa
// @Description Retorna uma tarefa específica
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da tarefa"
// @Success 200 {object} dto.TaskDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	task, err := h.taskService.GetTaskByID(c.Request().Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, appError.ErrTaskNotFound) {
			return response.NotFound(c, appError.ErrTaskNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return response.JSON(c, http.StatusOK, task)
}

// GetTasks lista todas as tarefas do usuário
// @Summary Listar tarefas
// @Description Retorna todas as tarefas do usuário com filtros opcionais
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param status_id query string false "Filtrar por status"
// @Param priority query string false "Filtrar por prioridade"
// @Param completed query bool false "Filtrar por completado"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {array} dto.TaskDTO
// @Router /api/v1/tasks [get]
func (h *TaskHandler) GetTasks(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	filters := &dto.TaskFilters{}

	if statusID := c.QueryParam("status_id"); statusID != "" {
		filters.StatusID = &statusID
	}

	if priority := c.QueryParam("priority"); priority != "" {
		filters.Priority = &priority
	}

	if completedStr := c.QueryParam("completed"); completedStr != "" {
		completed, err := strconv.ParseBool(completedStr)
		if err == nil {
			filters.Completed = &completed
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

	tasks, err := h.taskService.GetTasksByUserID(c.Request().Context(), userID, filters)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar tarefas", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, tasks)
}

// UpdateTask atualiza uma tarefa
// @Summary Atualizar tarefa
// @Description Atualiza os dados de uma tarefa
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da tarefa"
// @Param request body dto.UpdateTaskRequest true "Dados atualizados"
// @Success 200 {object} dto.TaskDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	task, err := h.taskService.UpdateTask(c.Request().Context(), userID, taskID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrTaskNotFound) {
			return response.NotFound(c, appError.ErrTaskNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao atualizar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, task)
}

// DeleteTask deleta uma tarefa
// @Summary Deletar tarefa
// @Description Remove uma tarefa
// @Tags tasks
// @Security BearerAuth
// @Param id path string true "ID da tarefa"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.taskService.DeleteTask(c.Request().Context(), userID, taskID); err != nil {
		if errors.Is(err, appError.ErrTaskNotFound) {
			return response.NotFound(c, appError.ErrTaskNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro ao atualizar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

// ReorderTask reordena uma tarefa (drag-and-drop)
// @Summary Reordenar tarefa
// @Description Reordena uma tarefa usando drag-and-drop
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ReorderTasksRequest true "Dados de reordenação"
// @Success 200 {object} dto.TaskDTO
// @Router /api/v1/tasks/reorder [post]
func (h *TaskHandler) ReorderTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.ReorderTasksRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	task, err := h.taskService.ReorderTask(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao reordenar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, task)
}

// CompleteTask marca uma tarefa como completa
// @Summary Completar tarefa
// @Description Marca uma tarefa como completa
// @Tags tasks
// @Security BearerAuth
// @Param id path string true "ID da tarefa"
// @Success 200 {object} dto.TaskDTO
// @Router /api/v1/tasks/{id}/complete [post]
func (h *TaskHandler) CompleteTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	task, err := h.taskService.CompleteTask(c.Request().Context(), userID, taskID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao completar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, task)
}

// UncompleteTask marca uma tarefa como incompleta
// @Summary Descompletar tarefa
// @Description Marca uma tarefa como incompleta
// @Tags tasks
// @Security BearerAuth
// @Param id path string true "ID da tarefa"
// @Success 200 {object} dto.TaskDTO
// @Router /api/v1/tasks/{id}/uncomplete [post]
func (h *TaskHandler) UncompleteTask(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	task, err := h.taskService.UncompleteTask(c.Request().Context(), userID, taskID)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao descompletar tarefa", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, task)
}

