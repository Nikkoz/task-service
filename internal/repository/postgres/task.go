package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/repository"
)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, t task.Task) (task.Task, error) {
	const q = `
		INSERT INTO tasks (title, description, status, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, status, due_date, created_at, updated_at
	`

	var out task.Task
	var status string

	err := r.db.QueryRow(ctx, q,
		t.Title,
		t.Description,
		string(t.Status),
		t.DueDate,
	).Scan(
		&out.ID,
		&out.Title,
		&out.Description,
		&status,
		&out.DueDate,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return task.Task{}, err
	}

	out.Status = task.Status(status)
	return out, nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id uint64) (task.Task, error) {
	const q = `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var out task.Task
	var status string

	err := r.db.QueryRow(ctx, q, id).Scan(
		&out.ID,
		&out.Title,
		&out.Description,
		&status,
		&out.DueDate,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return task.Task{}, repository.ErrNotFound
		}
		return task.Task{}, err
	}

	out.Status = task.Status(status)
	return out, nil
}

func (r *TaskRepo) List(ctx context.Context, limit, offset uint64) ([]task.Task, error) {
	const q = `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]task.Task, 0, limit)
	for rows.Next() {
		var t task.Task
		var status string

		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&status,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		t.Status = task.Status(status)
		out = append(out, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *TaskRepo) Update(ctx context.Context, t task.Task) (task.Task, error) {
	const q = `
		UPDATE tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    due_date = $5
		WHERE id = $1
		RETURNING id, title, description, status, due_date, created_at, updated_at
	`

	var out task.Task
	var status string

	err := r.db.QueryRow(ctx, q,
		t.ID,
		t.Title,
		t.Description,
		string(t.Status),
		t.DueDate,
	).Scan(
		&out.ID,
		&out.Title,
		&out.Description,
		&status,
		&out.DueDate,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return task.Task{}, repository.ErrNotFound
		}
		return task.Task{}, err
	}

	out.Status = task.Status(status)
	return out, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id uint64) error {
	const q = `DELETE FROM tasks WHERE id = $1`

	ct, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
