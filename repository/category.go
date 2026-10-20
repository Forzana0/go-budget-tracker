package repository

import (
	"context"
	"errors"
	"thirdProgram/models"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	GetAll(ctx context.Context, userID int) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category, userID int) (*models.Category, error)
	Delete(ctx context.Context, userID int, categoryID int) error
	Update(ctx context.Context, userID int, categoryID int, category *models.Category) (*models.Category, error)
}
