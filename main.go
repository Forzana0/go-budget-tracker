package main

import (
	"context"
	"log"
	"net/http"
	"thirdProgram/config"
	"thirdProgram/database"
	"thirdProgram/handlers"
	"thirdProgram/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	// 1. Завантажуємо конфігурацію застосунку
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	ctx := context.Background()

	// 2. Використовуємо DatabaseURL з конфігурації замість хардкоду
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewPostgresUserRepository(pool)
	categoryRepo := repository.NewPostgresCategoryRepository(pool)
	expenseRepo := repository.NewPostgresExpenseRepository(pool)

	h := handlers.New(userRepo, categoryRepo, expenseRepo)

	r := chi.NewRouter()

	// User_Handlers
	r.Get("/users", h.GetUsers)
	r.Post("/users", h.CreateUser)
	r.Delete("/users/{user_id}", h.DeleteUser)
	r.Put("/users/{user_id}", h.UpdateUser)

	// Category_handlers
	r.Get("/users/{user_id}/categories", h.GetCategories)
	r.Post("/users/{user_id}/categories", h.CreateCategory)
	r.Delete("/users/{user_id}/categories/{category_id}", h.DeleteCategory)
	r.Put("/users/{user_id}/categories/{category_id}", h.UpdateCategory)

	// Expense_handlers
	r.Post("/users/{user_id}/expenses", h.AddExpense)
	r.Get("/users/{user_id}/expenses", h.GetExpenses)
	r.Delete("/users/{user_id}/expenses/{expense_id}", h.DeleteExpense)
	r.Put("/users/{user_id}/expenses/{expense_id}", h.UpdateExpense)

	// 3. Запускаємо сервер на порті з конфігурації
	log.Printf("Server is starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, r); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
