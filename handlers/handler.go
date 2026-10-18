package handlers

import (
	"thirdProgram/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
	pool         *pgxpool.Pool // Залишаємо поки що для категорій та витрат
}

func New(userRepo repository.UserRepository, categoryRepo repository.CategoryRepository, pool *pgxpool.Pool) *Handler {
	return &Handler{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		pool:         pool,
	}
}
