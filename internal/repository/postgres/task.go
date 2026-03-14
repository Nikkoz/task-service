package postgres

import (
	"errors"
	"time"

	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/Nikkoz/task-service/pkg/logger"
	"github.com/jackc/pgx/v5"
)

var _ service.TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	db DBTX
}

func NewTaskRepo(db DBTX) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, t task.Task) (task.Task, error) {
	const q = `
		INSERT INTO tasks (user_id, title, description, status, due_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, description, status, due_date, created_at, updated_at
	`

	var (
		out    task.Task
		status string
		due    *time.Time
	)

	err := r.db.QueryRow(ctx, q,
		t.UserID,
		t.Title,
		t.Description,
		string(t.Status),
		t.DueDate.DateTime(),
	).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Description,
		&status,
		&due,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return task.Task{}, logger.ErrorWithContext(ctx, err)
	}

	out.Status = task.Status(status)
	if due != nil {
		out.DueDate, _ = task.NewDueDate(*due)
	}

	return out, nil
}

func (r *TaskRepo) Update(ctx context.Context, t task.Task) (task.Task, error) {
	const q = `
		UPDATE tasks
		SET title = $3,
		    description = $4,
		    status = $5,
		    due_date = $6
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, title, description, status, due_date, created_at, updated_at
	`
	var (
		out    task.Task
		status string
		due    *time.Time
	)

	err := r.db.QueryRow(ctx, q,
		t.ID,
		t.UserID,
		t.Title,
		t.Description,
		string(t.Status),
		t.DueDate.DateTime(),
	).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Description,
		&status,
		&due,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return task.Task{}, repository.ErrNotFound
		}

		return task.Task{}, logger.ErrorWithContext(ctx, err)
	}

	out.Status = task.Status(status)
	if due != nil {
		out.DueDate, _ = task.NewDueDate(*due)
	}

	return out, nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id, userID uint64) (task.Task, error) {
	const q = `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`

	var (
		out    task.Task
		status string
		due    *time.Time
	)

	err := r.db.QueryRow(ctx, q, id, userID).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Description,
		&status,
		&due,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return task.Task{}, repository.ErrNotFound
		}

		return task.Task{}, logger.ErrorWithContext(ctx, err)
	}

	out.Status = task.Status(status)
	if due != nil {
		out.DueDate, _ = task.NewDueDate(*due)
	}

	return out, nil
}

func (r *TaskRepo) List(ctx context.Context, userID uint64, limit, offset uint64) ([]task.Task, error) {
	const q = `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, logger.ErrorWithContext(ctx, err)
	}
	defer rows.Close()

	out := make([]task.Task, 0, limit)
	for rows.Next() {
		var (
			t      task.Task
			status string
			due    *time.Time
		)

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&status,
			&due,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, logger.ErrorWithContext(ctx, err)
		}

		t.Status = task.Status(status)
		if due != nil {
			t.DueDate, _ = task.NewDueDate(*due)
		}

		out = append(out, t)
	}

	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWithContext(ctx, err)
	}

	return out, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id, userID uint64) error {
	const q = `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	ct, err := r.db.Exec(ctx, q, id, userID)
	if err != nil {
		return logger.ErrorWithContext(ctx, err)
	}

	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
