package project

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

func (r *PostgresRepository) Create(ctx context.Context, project Project) (Project, error) {
	const query = `
		INSERT INTO projects (
			name,
			description,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		project.Name,
		project.Description,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(&project.ID)
	if err != nil {
		return Project{}, fmt.Errorf(
			"create project: %w",
			err,
		)
	}
	return project, nil
}

func (r *PostgresRepository) FindAll(
	ctx context.Context,
) ([]Project, error) {
	const query = `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"find all projects: %w",
			err,
		)
	}
	defer rows.Close()

	projects := make([]Project, 0)

	for rows.Next() {
		var project Project

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan project: %w",
				err,
			)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate projects: %w",
			err,
		)
	}

	return projects, nil
}

func (r *PostgresRepository) FindByID(
	ctx context.Context,
	id int64,
) (Project, error) {
	const query = `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1
	`

	var project Project

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Project{}, ErrNotFound
	}

	if err != nil {
		return Project{}, fmt.Errorf(
			"find project by ID: %w",
			err,
		)
	}

	return project, nil
}

func (r *PostgresRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM projects
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrHasTasks
		}
		return fmt.Errorf(
			"delete project: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get deleted project count: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
