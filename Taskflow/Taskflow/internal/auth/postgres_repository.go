package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, user User) (User, error) {
	const query = `
		INSERT INTO users (
			name,
			email,
			password_hash,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.Is(err, pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailAlreadyExists
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *PostgresRepository) FindByID(
	ctx context.Context,
	id int64,
) (User, error) {
	const query = `
		SELECT
			id,
			name,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf(
			"find user by ID: %w",
			err,
		)
	}

	return user, nil
}

func (r *PostgresRepository) FindByEmail(
	ctx context.Context,
	email string,
) (User, error) {
	const query = `
		SELECT
			id,
			name,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`

	var user User

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf(
			"find user by email: %w",
			err,
		)
	}

	return user, nil
}
