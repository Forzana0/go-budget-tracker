package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"thirdProgram/models"
	"thirdProgram/repository"
)

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	expanse := &models.Expense{}
	if err := json.NewDecoder(r.Body).Decode(expanse); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	created, err := h.expenseRepo.Create(ctx, expanse, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotEnoughMoney) {
			http.Error(w, "Not enough money", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error while creating expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()

	expenses, err := h.expenseRepo.GetAll(ctx, userID)
	if err != nil {
		http.Error(w, "Error getting expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		http.Error(w, " Error encoding json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	expenseID, ok := parseExpenseID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()

	err := h.expenseRepo.Delete(ctx, userID, expenseID)
	if err != nil {
		http.Error(w, "Error while deleting expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	ctx := r.Context()

	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	expenseID, ok := parseExpenseID(w, r)
	if !ok {
		return
	}

	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updated, err := h.expenseRepo.Update(ctx, userID, expenseID, &expense)
	if err != nil {
		if errors.Is(err, repository.ErrNotEnoughMoney) {
			http.Error(w, "Not enough money", http.StatusForbidden)
			return
		}
		http.Error(w, "Error while updating expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}
