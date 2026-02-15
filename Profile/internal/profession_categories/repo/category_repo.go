package repo

import (
	"context"
	"database/sql"

	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/models"
)

type CategoryRepo interface {
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
	GetCategoryByID(ctx context.Context, id int16) (*models.Category, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetCategoriesByParent(ctx context.Context, parentID int16) ([]models.Category, error)
	GetParentOfCategory(ctx context.Context, id int16) (*int16, error)
}

type PostgresCategoryRepo struct {
	db *sql.DB
}

func NewPostgresCategoryRepo(db *sql.DB) *PostgresCategoryRepo {
	return &PostgresCategoryRepo{db: db}
}
