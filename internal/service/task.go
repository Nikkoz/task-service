package service

import (
	"time"

	domain "github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/pkg/context"
)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

type CreateTaskInput struct {
	UserID      uint64
	Title       string
	Description string
	Status      string
	DueDate     *time.Time
}

type UpdateTaskInput struct {
	Title       string
	Description string
	Status      string
	DueDate     *time.Time
}

func (s *TaskService) CreateTask(ctx context.Context, in CreateTaskInput) (domain.Task, error) {
	taskTitle, err := domain.NewTitle(in.Title)
	if err != nil {
		return domain.Task{}, err
	}

	taskDescription, err := domain.NewDescription(in.Description)
	if err != nil {
		return domain.Task{}, err
	}

	taskStatus, err := domain.NewStatus(in.Status)
	if err != nil {
		return domain.Task{}, err
	}

	t := domain.Task{
		UserID:      in.UserID,
		Title:       *taskTitle,
		Description: *taskDescription,
		Status:      *taskStatus,
	}

	if in.DueDate == nil {
		return s.repo.Create(ctx, t)
	}

	dueDate, err := domain.NewDueDate(*in.DueDate)
	if err != nil {
		return domain.Task{}, err
	}

	t.DueDate = dueDate

	return s.repo.Create(ctx, t)
}

func (s *TaskService) UpdateTask(ctx context.Context, id, userID uint64, in UpdateTaskInput) (domain.Task, error) {
	if userID <= 0 || id <= 0 {
		return domain.Task{}, ErrValidation
	}

	taskTitle, err := domain.NewTitle(in.Title)
	if err != nil {
		return domain.Task{}, err
	}

	taskDescription, err := domain.NewDescription(in.Description)
	if err != nil {
		return domain.Task{}, err
	}

	taskStatus, err := domain.NewStatus(in.Status)
	if err != nil {
		return domain.Task{}, err
	}

	t := domain.Task{
		ID:          id,
		UserID:      userID,
		Title:       *taskTitle,
		Description: *taskDescription,
		Status:      *taskStatus,
	}

	if in.DueDate == nil {
		return s.repo.Update(ctx, t)
	}

	dueDate, err := domain.NewDueDate(*in.DueDate)
	if err != nil {
		return domain.Task{}, err
	}

	t.DueDate = dueDate

	return s.repo.Update(ctx, t)
}

func (s *TaskService) GetTask(ctx context.Context, id, userID uint64) (domain.Task, error) {
	if userID <= 0 || id <= 0 {
		return domain.Task{}, ErrValidation
	}

	return s.repo.GetByID(ctx, id, userID)
}

func (s *TaskService) ListTasks(ctx context.Context, userID uint64, limit, offset uint64) ([]domain.Task, error) {
	if userID <= 0 {
		return nil, ErrValidation
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.repo.List(ctx, userID, limit, offset)
}

func (s *TaskService) DeleteTask(ctx context.Context, id, userID uint64) error {
	if userID <= 0 || id <= 0 {
		return ErrValidation
	}

	return s.repo.Delete(ctx, id, userID)
}
