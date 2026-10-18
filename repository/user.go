package repository

import (
	"context"
	"thirdProgram/models"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int) error
}
