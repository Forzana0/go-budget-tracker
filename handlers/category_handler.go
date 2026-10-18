package handlers

import (
	"encoding/json"
	"net/http"
	"thirdProgram/models"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()

	categories, err := h.categoryRepo.GetAll(ctx, userID)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Error while encoding json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()
	category := &models.Category{}

	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	created, err := h.categoryRepo.Create(ctx, category, userID)
	if err != nil {
		http.Error(w, "Error while creatin category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, ok := parseCategoryID(w, r)
	if !ok {
		return
	}
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()

	var dummy int
	err := h.pool.QueryRow(ctx, "SELECT id FROM categories WHERE id = $1 AND user_id = $2", categoryID, userID).Scan(&dummy)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Category not found or it is a global category you can not delete", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		http.Error(w, "Error starting transaction for deleting category", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE expenses SET category_id = $1 WHERE category_id =$2 AND user_id = $3", 6, categoryID, userID)
	if err != nil {
		http.Error(w, "Error while changing category to undefined", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM categories WHERE id = $1 AND user_id = $2", categoryID, userID)
	if err != nil {
		http.Error(w, "Error while deleting category", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Error commiting transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
