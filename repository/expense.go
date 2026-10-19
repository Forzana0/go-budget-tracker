package repository

import (
	"context"
	"errors"
	"thirdProgram/models"
)

var ErrNotEnoughMoney = errors.New("not enough money")

type ExpenseRepository interface {
	GetAll(ctx context.Context, userID int) ([]models.Expense, error)
	Create(ctx context.Context, expense *models.Expense, userID int) (*models.Expense, error)
	Delete(ctx context.Context, userID int, expenseID int) error
}
