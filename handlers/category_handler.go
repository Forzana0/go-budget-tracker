package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"thirdProgram/models"
	"thirdProgram/repository"

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

	err := h.categoryRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Category not found or it is a category you cannot delete", http.StatusForbidden)
			return
		}
		http.Error(w, "Error while deleting category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, ok := parseCategoryID(w, r)
	if !ok {
		return
	}

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	var category models.Category
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updated, err := h.categoryRepo.Update(ctx, userID, categoryID, &category)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error while updating category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)

}
