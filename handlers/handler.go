package handlers

import (
	"thirdProgram/repository"
)

type Handler struct {
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
	expenseRepo  repository.ExpenseRepository
}

func New(userRepo repository.UserRepository, categoryRepo repository.CategoryRepository, expenseRepo repository.ExpenseRepository) *Handler {
	return &Handler{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		expenseRepo:  expenseRepo,
	}
}
