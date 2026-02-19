package categoryService

import (
	"context"

	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/repo"
)

type CategoryService interface {
	CreateCategoryService(ctx context.Context, req models.CategoryDTO) (*models.Category, error)
	GetCategoriesService(ctx context.Context) ([]models.Category, error)
	GetCategoriesByParentService(ctx context.Context, parentID int16) ([]models.Category, error)
	GetCategoryByIDService(ctx context.Context, id int16) (*models.Category, error)
}

type ProfessionCategoryService struct {
	repo repo.CategoryRepo
}

func NewProfessionCategoryService(repo repo.CategoryRepo) *ProfessionCategoryService {
	return &ProfessionCategoryService{
		repo: repo,
	}
}

func (s *ProfessionCategoryService) CreateCategoryService(ctx context.Context, req models.CategoryDTO) (*models.Category, error) {
	category := &models.Category{
		Name: req.Name,
	}

	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	res, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProfessionCategoryService) GetCategoriesService(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *ProfessionCategoryService) GetCategoriesByParentService(ctx context.Context, parentID int16) ([]models.Category, error) {
	return s.repo.GetCategoriesByParent(ctx, parentID)
}

func (s *ProfessionCategoryService) GetCategoryByIDService(ctx context.Context, id int16) (*models.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}
