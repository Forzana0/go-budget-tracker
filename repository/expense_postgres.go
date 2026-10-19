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
