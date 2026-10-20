package repository

import (
	"context"
	"errors"
	"thirdProgram/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	rows, err := r.pool.Query(ctx, "SELECT id, username, balance FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Balance); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	var created models.User

	err := r.pool.QueryRow(
		ctx,
		"INSERT INTO users (username, balance) VALUES ($1, $2) RETURNING id, username, balance",
		user.Username,
		user.Balance,
	).Scan(&created.ID, &created.Username, &created.Balance)

	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, userID int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}

func (r *PostgresUserRepository) Update(ctx context.Context, userID int, user *models.User) (*models.User, error) {
	var updated models.User

	err := r.pool.QueryRow(
		ctx,
		"UPDATE users SET username = $1, balance = $2 WHERE id =$3 RETURNING id, username, balance",
		user.Username,
		user.Balance,
		userID,
	).Scan(&updated.ID, &updated.Username, &updated.Balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &updated, nil
}
