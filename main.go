package main

import (
	"context"
	"net/http"
	"thirdProgram/database"
	"thirdProgram/handlers"
	"thirdProgram/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	pool, err := database.Connect(ctx, "postgres://postgres:@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	userRepo := repository.NewPostgresUserRepository(pool)
	categoryRepo := repository.NewPostgresCategoryRepository(pool)

	h := handlers.New(userRepo, categoryRepo, pool)

	r := chi.NewRouter()

	//User_Handlers
	r.Get("/users", h.GetUsers)
	r.Post("/users", h.CreateUser)
	r.Delete("/users", h.DeleteUser)

	//Category_handlers
	r.Get("/users/{user_id}/categories", h.GetCategories)
	r.Post("/users/{user_id}/categories", h.CreateCategory)
	r.Delete("/users/{user_id}/categories/{category_id}", h.DeleteCategory)

	//Expense_handler
	r.Post("/users/{user_id}/expenses", h.AddExpense)
	r.Get("/users/{user_id}/expenses", h.GetExpenses)
	r.Delete("/users/{user_id}/expenses/{expense_id}", h.DeleteExpense)

	http.ListenAndServe(":9091", r)
}
