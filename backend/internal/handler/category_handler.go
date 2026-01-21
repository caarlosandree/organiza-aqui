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

type CategoryHandler struct {
	categoryService service.CategoryService
	validator       *validator.Validate
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validator:       validator.New(),
	}
}

// CreateCategory cria uma nova categoria
// @Summary Criar categoria
// @Description Cria uma nova categoria financeira
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoryDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req dto.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	category, err := h.categoryService.CreateCategory(c.Request().Context(), userID, &req)
	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao criar categoria", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusCreated, category)
}

// GetCategory obtém uma categoria por ID
// @Summary Obter categoria
// @Description Retorna uma categoria específica
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoryDTO
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	categoryID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	category, err := h.categoryService.GetCategoryByID(c.Request().Context(), userID, categoryID)
	if err != nil {
		if errors.Is(err, appError.ErrCategoryNotFound) {
			return response.NotFound(c, appError.ErrCategoryNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, category)
}

// GetCategories lista todas as categorias do usuário
// @Summary Listar categorias
// @Description Retorna todas as categorias do usuário autenticado
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filtrar por tipo (income, expense)"
// @Param tree query bool false "Retornar como árvore hierárquica"
// @Success 200 {array} dto.CategoryDTO
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	categoryType := c.QueryParam("type")
	tree := c.QueryParam("tree") == "true"

	var categories []*dto.CategoryDTO

	if tree {
		categories, err = h.categoryService.BuildCategoryTree(c.Request().Context(), userID)
	} else if categoryType != "" {
		categories, err = h.categoryService.GetCategoriesByUserIDAndType(c.Request().Context(), userID, categoryType)
	} else {
		categories, err = h.categoryService.GetCategoriesByUserID(c.Request().Context(), userID)
	}

	if err != nil {
		logger.ErrorWithContext(c.Request().Context(), "Erro ao buscar categorias", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, categories)
}

// UpdateCategory atualiza uma categoria
// @Summary Atualizar categoria
// @Description Atualiza os dados de uma categoria
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da categoria"
// @Param request body dto.UpdateCategoryRequest true "Dados atualizados da categoria"
// @Success 200 {object} dto.CategoryDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	categoryID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.BadRequest(c, appError.ErrInvalidInput)
	}

	category, err := h.categoryService.UpdateCategory(c.Request().Context(), userID, categoryID, &req)
	if err != nil {
		if errors.Is(err, appError.ErrCategoryNotFound) {
			return response.NotFound(c, appError.ErrCategoryNotFound)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

		return response.JSON(c, http.StatusOK, category)
}

// DeleteCategory deleta uma categoria
// @Summary Deletar categoria
// @Description Remove uma categoria do sistema
// @Tags categories
// @Security BearerAuth
// @Param id path string true "ID da categoria"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	categoryID, err := parseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.categoryService.DeleteCategory(c.Request().Context(), userID, categoryID); err != nil {
		if errors.Is(err, appError.ErrCategoryNotFound) {
			return response.NotFound(c, appError.ErrCategoryNotFound)
		}
		if errors.Is(err, appError.ErrHasSubcategories) {
			return response.BadRequest(c, appError.ErrHasSubcategories)
		}
		logger.ErrorWithContext(c.Request().Context(), "Erro interno", err)
		return response.InternalServerError(c, appError.ErrInternalServer)
	}

	return c.NoContent(http.StatusNoContent)
}

