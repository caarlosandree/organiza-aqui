package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/organiza-aqui/backend/internal/dto"
	appError "github.com/organiza-aqui/backend/internal/error"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/organiza-aqui/backend/internal/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, userID uuid.UUID, req *dto.CreateCategoryRequest) (*dto.CategoryDTO, error)
	GetCategoryByID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (*dto.CategoryDTO, error)
	GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.CategoryDTO, error)
	GetCategoriesByUserIDAndType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*dto.CategoryDTO, error)
	UpdateCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, req *dto.UpdateCategoryRequest) (*dto.CategoryDTO, error)
	DeleteCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) error
	BuildCategoryTree(ctx context.Context, userID uuid.UUID) ([]*dto.CategoryDTO, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, userID uuid.UUID, req *dto.CreateCategoryRequest) (*dto.CategoryDTO, error) {
	var parentID *uuid.UUID
	var path string

	if req.ParentID != nil {
		parsedParentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent_id inválido: %w", err)
		}

		parent, err := s.categoryRepo.FindByID(ctx, parsedParentID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar categoria pai: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("categoria pai %s: %w", parsedParentID, appError.ErrCategoryNotFound)
		}
		if parent.UserID != userID {
			return nil, appError.ErrUnauthorizedAccess
		}

		parentID = &parsedParentID
		// Calcular path: pegar o último número do path do pai e incrementar
		pathParts := strings.Split(parent.Path, ".")
		lastNum, _ := strconv.Atoi(pathParts[len(pathParts)-1])
		path = fmt.Sprintf("%s.%d", parent.Path, lastNum+1)
	} else {
		// Categoria raiz - buscar o maior número de path para este usuário
		categories, err := s.categoryRepo.FindByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
		}

		maxNum := 0
		for _, cat := range categories {
			if cat.ParentID == nil {
				pathParts := strings.Split(cat.Path, ".")
				if len(pathParts) > 0 {
					if num, err := strconv.Atoi(pathParts[0]); err == nil && num > maxNum {
						maxNum = num
					}
				}
			}
		}
		path = strconv.Itoa(maxNum + 1)
	}

	now := time.Now()
	category := &model.Category{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		ParentID:  parentID,
		Path:      path,
		Type:      req.Type,
		Color:     req.Color,
		CreatedAt: now,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	return s.modelToDTO(category), nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (*dto.CategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("categoria %s: %w", categoryID, appError.ErrCategoryNotFound)
	}
	if category.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	return s.modelToDTO(category), nil
}

func (s *categoryService) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.CategoryDTO, error) {
	categories, err := s.categoryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}

	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = s.modelToDTO(category)
	}

	return dtos, nil
}

func (s *categoryService) GetCategoriesByUserIDAndType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*dto.CategoryDTO, error) {
	categories, err := s.categoryRepo.FindByUserIDAndType(ctx, userID, categoryType)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}

	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = s.modelToDTO(category)
	}

	return dtos, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, req *dto.UpdateCategoryRequest) (*dto.CategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("categoria %s: %w", categoryID, appError.ErrCategoryNotFound)
	}
	if category.UserID != userID {
		return nil, appError.ErrUnauthorizedAccess
	}

	var parentID *uuid.UUID
	if req.ParentID != nil {
		parsedParentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent_id inválido: %w", err)
		}
		parentID = &parsedParentID
	}

	category.Name = req.Name
	category.ParentID = parentID
	category.Type = req.Type
	category.Color = req.Color

	// Recalcular path se parent mudou
	if parentID != nil {
		parent, err := s.categoryRepo.FindByID(ctx, *parentID)
		if err == nil && parent != nil {
			pathParts := strings.Split(parent.Path, ".")
			lastNum, _ := strconv.Atoi(pathParts[len(pathParts)-1])
			category.Path = fmt.Sprintf("%s.%d", parent.Path, lastNum+1)
		}
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	return s.modelToDTO(category), nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) error {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("erro ao buscar categoria: %w", err)
	}
	if category == nil {
		return fmt.Errorf("categoria %s: %w", categoryID, appError.ErrCategoryNotFound)
	}
	if category.UserID != userID {
		return appError.ErrUnauthorizedAccess
	}

	// Verificar se tem filhos
	children, err := s.categoryRepo.FindChildren(ctx, categoryID)
	if err == nil && len(children) > 0 {
		return appError.ErrHasSubcategories
	}

	if err := s.categoryRepo.Delete(ctx, categoryID); err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	return nil
}

func (s *categoryService) BuildCategoryTree(ctx context.Context, userID uuid.UUID) ([]*dto.CategoryDTO, error) {
	categories, err := s.categoryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}

	// Criar mapa de categorias por ID
	categoryMap := make(map[uuid.UUID]*dto.CategoryDTO)
	for _, cat := range categories {
		catDTO := s.modelToDTO(cat)
		catDTO.Children = []dto.CategoryDTO{}
		categoryMap[cat.ID] = catDTO
	}

	// Construir árvore
	var roots []*dto.CategoryDTO
	for _, cat := range categories {
		catDTO := categoryMap[cat.ID]
		if cat.ParentID == nil {
			roots = append(roots, catDTO)
		} else {
			parent := categoryMap[*cat.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, *catDTO)
			}
		}
	}

	return roots, nil
}

func (s *categoryService) modelToDTO(category *model.Category) *dto.CategoryDTO {
	dto := &dto.CategoryDTO{
		ID:        category.ID.String(),
		UserID:    category.UserID.String(),
		Name:      category.Name,
		Path:      category.Path,
		Type:      category.Type,
		Color:     category.Color,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		Children:  []dto.CategoryDTO{},
	}

	if category.ParentID != nil {
		parentIDStr := category.ParentID.String()
		dto.ParentID = &parentIDStr
	}

	return dto
}
