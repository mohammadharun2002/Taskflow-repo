package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, task Task) (Task, error) {
	const query = `
		INSERT INTO tasks(
			project_id,
			name,
			description,
			status,
			priority,
			due_date,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		task.ProjectID,
		task.Name,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (r *PostgresRepository) FindByProjectID(
	ctx context.Context,
	projectID int64,
) ([]Task, error) {
	const query = `
		SELECT
			id,
			project_id,
			name,
			description,
			status,
			priority,
			due_date,
			created_at,
			updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"find tasks by project ID: %w",
			err,
		)
	}
	defer rows.Close()

	tasks := make([]Task, 0)

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan task: %w",
				err,
			)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate tasks: %w",
			err,
		)
	}

	return tasks, nil
}

func (r *PostgresRepository) FindByID(
	ctx context.Context,
	id int64,
) (Task, error) {
	const query = `
		SELECT
			id,
			project_id,
			name,
			description,
			status,
			priority,
			due_date,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
	`

	var task Task

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}

	if err != nil {
		return Task{}, fmt.Errorf(
			"find task by ID: %w",
			err,
		)
	}

	return task, nil
}

func (r *PostgresRepository) Update(
	ctx context.Context,
	task Task,
) (Task, error) {
	const query = `
		UPDATE tasks
		SET
			name = $1,
			description = $2,
			status = $3,
			priority = $4,
			due_date = $5,
			updated_at = $6
		WHERE id = $7
		RETURNING
			id,
			project_id,
			name,
			description,
			status,
			priority,
			due_date,
			created_at,
			updated_at
	`

	var updatedTask Task

	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
	).Scan(
		&updatedTask.ID,
		&updatedTask.ProjectID,
		&updatedTask.Name,
		&updatedTask.Description,
		&updatedTask.Status,
		&updatedTask.Priority,
		&updatedTask.DueDate,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}

	if err != nil {
		return Task{}, fmt.Errorf(
			"update task: %w",
			err,
		)
	}

	return updatedTask, nil
}

func (r *PostgresRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf(
			"delete task: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get deleted task count: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
