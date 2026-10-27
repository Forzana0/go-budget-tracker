package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"thirdProgram/models"
	"thirdProgram/repository"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.userRepo.GetAll(ctx)

	if err != nil {
		http.Error(w, "Error request to DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	createdUser, err := h.userRepo.Create(ctx, &user)
	if err != nil {
		http.Error(w, "Error while adding new user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	err := h.userRepo.Delete(ctx, userID)
	if err != nil {
		http.Error(w, "Error while deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	updated, err := h.userRepo.Update(ctx, userID, &user)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error while updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}
