package repository

import (
	"context"
	"thirdProgram/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresExpenseRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresExpenseRepository(pool *pgxpool.Pool) *PostgresExpenseRepository {
	return &PostgresExpenseRepository{
		pool: pool,
	}
}

func (r *PostgresExpenseRepository) GetAll(ctx context.Context, userID int) ([]models.Expense, error) {

	var expenses []models.Expense
	rows, err := r.pool.Query(ctx, "SELECT id, name, description, price, created_at, user_id, category_id FROM expenses WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Price, &e.CreatedAt, &e.UserID, &e.CategoryID); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *PostgresExpenseRepository) Create(ctx context.Context, expense *models.Expense, userID int) (*models.Expense, error) {
	var currentBalance float64
	var created models.Expense
	err := r.pool.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
	if err != nil {
		return nil, err
	}

	if expense.Price > currentBalance {
		return nil, ErrNotEnoughMoney
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", expense.Price, userID)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, "INSERT INTO expenses(name, description, price, user_id, category_id) VALUES($1, $2, $3, $4, $5) RETURNING id, name, description, price, created_at, user_id, category_id", expense.Name, expense.Description, expense.Price, userID, expense.CategoryID).Scan(&created.ID, &created.Name, &created.Description, &created.Price, &created.CreatedAt, &created.UserID, &created.CategoryID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresExpenseRepository) Delete(ctx context.Context, userID int, expenseID int) error {
	var price float64
	err := r.pool.QueryRow(ctx, "SELECT price FROM expenses WHERE id = $1", expenseID).Scan(&price)
	if err != nil {
		return err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", price, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM expenses WHERE id = $1 AND user_id = $2", expenseID, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresExpenseRepository) Update(ctx context.Context, userID int, expenseID int, expense *models.Expense) (*models.Expense, error) {
	var updated models.Expense
	var oldPrice float64
	var balance float64

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "SELECT price FROM expenses WHERE id = $1 AND user_id = $2", expenseID, userID).Scan(&oldPrice)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, "UPDATE expenses SET name = $1, description = $2, price = $3, category_id = $4 WHERE id = $5 AND user_id = $6 RETURNING id, name, description, price, created_at, user_id, category_id", expense.Name, expense.Description, expense.Price, expense.CategoryID, expenseID, userID).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.Price, &updated.CreatedAt, &updated.UserID, &updated.CategoryID)
	if err != nil {
		return nil, err
	}

	priceDiff := oldPrice - updated.Price
	newBalance := balance + priceDiff

	if newBalance >= 0 {
		_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE  id = $2", priceDiff, userID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrNotEnoughMoney
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &updated, nil
}
