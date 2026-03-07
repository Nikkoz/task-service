package task

import (
	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/pkg/context"
)

//go:generate mockery --name Service --output ./mocks --outpkg mocks
type Service interface {
	CreateTask(ctx context.Context, in service.CreateTaskInput) (task.Task, error)
	GetTask(ctx context.Context, id uint64) (task.Task, error)
	ListTasks(ctx context.Context, limit, offset uint64) ([]task.Task, error)
	UpdateTask(ctx context.Context, id uint64, in service.UpdateTaskInput) (task.Task, error)
	DeleteTask(ctx context.Context, id uint64) error
}
