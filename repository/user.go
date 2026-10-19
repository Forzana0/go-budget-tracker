package repository

import (
	"context"
	"errors"
	"thirdProgram/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, userID int, user *models.User) (*models.User, error)
}
