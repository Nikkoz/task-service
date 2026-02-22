package service

import (
	"context"

	"github.com/Nikkoz/task-service/internal/domain/task"
)

type (
	Task interface {
		Create(ctx context.Context, t task.Task) (task.Task, error)
		Update(ctx context.Context, t task.Task) (task.Task, error)
		Delete(ctx context.Context, id uint64) error

		TaskReader
	}

	TaskReader interface {
		GetByID(ctx context.Context, id uint64) (task.Task, error)
		List(ctx context.Context, limit, offset uint64) ([]task.Task, error)
	}
)
