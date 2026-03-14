package service

import (
	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/pkg/context"
)

//go:generate mockery --name TaskRepository --output ./mocks --outpkg mocks
type TaskRepository interface {
	Create(ctx context.Context, t task.Task) (task.Task, error)
	Update(ctx context.Context, t task.Task) (task.Task, error)
	Delete(ctx context.Context, id, userID uint64) error

	TaskReaderRepository
}

type TaskReaderRepository interface {
	GetByID(ctx context.Context, id, userID uint64) (task.Task, error)
	List(ctx context.Context, userID uint64, limit, offset uint64) ([]task.Task, error)
}
