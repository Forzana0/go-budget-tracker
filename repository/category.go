package repository

import (
	"context"
	"thirdProgram/models"
)

type CategoryRepository interface {
	GetAll(ctx context.Context, userID int) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category, userID int) (*models.Category, error)
	Delete(ctx context.Context, userID int, categoryID int) error
}
