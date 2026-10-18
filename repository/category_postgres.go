package repository

import (
	"context"
	"thirdProgram/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCategoryRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCategoryRepository(pool *pgxpool.Pool) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{
		pool: pool,
	}
}

func (r *PostgresCategoryRepository) GetAll(ctx context.Context, userID int) ([]models.Category, error) {

	var categories []models.Category
	rows, err := r.pool.Query(ctx, "SELECT id, name, user_id FROM categories WHERE user_id is NULL OR user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.UserID); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) Create(ctx context.Context, category *models.Category, userID int) (*models.Category, error) {

	var created models.Category

	err := r.pool.QueryRow(
		ctx,
		"INSERT INTO categories (name, user_id) VALUES ($1, $2) RETURNING id, name, user_id",
		category.Name,
		userID,
	).Scan(&created.ID, &created.Name, &created.UserID)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
