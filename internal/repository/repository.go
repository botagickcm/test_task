package repository

import (
	"context"

	"test_task/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, login, password, last_login)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`

	err := r.pool.QueryRow(ctx, query,
		user.FirstName,
		user.LastName,
		user.Login,
		user.Password,
		user.LastLogin,
	).Scan(&user.ID)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT user_id, first_name, last_name, login, password, last_login
		FROM users
		WHERE user_id = $1
	`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Login,
		&user.Password,
		&user.LastLogin,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT user_id, first_name, last_name, login, password, last_login
		FROM users
		WHERE login = $1
	`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Login,
		&user.Password,
		&user.LastLogin,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT user_id, first_name, last_name, login, password, last_login
		FROM users
		ORDER BY user_id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Login,
			&user.Password,
			&user.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, login = $3, password = $4, last_login = $5
		WHERE user_id = $6
	`

	_, err := r.pool.Exec(ctx, query,
		user.FirstName,
		user.LastName,
		user.Login,
		user.Password,
		user.LastLogin,
		user.ID,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE user_id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
