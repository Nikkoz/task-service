package service

import (
	"testing"
	"time"

	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/service/mocks"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask_Success(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	service := NewTaskService(repo)

	title, _ := task.NewTitle("Buy milk")
	descr, _ := task.NewDescription("2 liters")
	status := task.StatusPlanned

	in := CreateTaskInput{
		Title:       title.String(),
		Description: descr.String(),
		Status:      status.String(),
		DueDate:     nil,
	}

	repo.
		On(
			"Create",
			mock.Anything,
			mock.MatchedBy(func(tsk task.Task) bool {
				return tsk.Title.String() == in.Title &&
					tsk.Description.String() == in.Description &&
					tsk.Status.String() == in.Status &&
					tsk.DueDate == nil
			}),
		).
		Return(task.Task{
			ID:          1,
			Title:       *title,
			Description: *descr,
			Status:      status,
			DueDate:     nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil).
		Once()

	assertion := assert.New(t)
	out, err := service.CreateTask(context.Empty(), in)

	assertion.NoError(err)
	assertion.Equal(uint64(1), out.ID)
	assertion.Equal(in.Title, out.Title.String())
	assertion.Equal(in.Description, out.Description.String())
	assertion.Equal(in.Status, out.Status.String())

	repo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_Success(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	service := NewTaskService(repo)

	title, _ := task.NewTitle("Updated title")
	descr, _ := task.NewDescription("2 liters")
	due, _ := task.NewDueDate(time.Now().UTC().Add(24 * time.Hour))
	status := task.StatusPlanned
	now := time.Now()

	dueDateTime := due.DateTime()
	in := UpdateTaskInput{
		Title:       title.String(),
		Description: descr.String(),
		Status:      status.String(),
		DueDate:     dueDateTime,
	}

	repo.
		On(
			"Update",
			mock.Anything,
			mock.MatchedBy(func(task task.Task) bool {
				if task.ID != 1 {
					return false
				}

				if task.Title.String() != in.Title {
					return false
				}

				if task.Description.String() != in.Description {
					return false
				}

				if task.Status.String() != in.Status {
					return false
				}

				if task.DueDate == nil {
					return false
				}

				return true
			}),
		).
		Return(task.Task{
			ID:          1,
			Title:       *title,
			Description: *descr,
			Status:      status,
			DueDate:     due,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil).
		Once()

	assertion := assert.New(t)
	out, err := service.UpdateTask(context.Empty(), 1, in)

	assertion.NoError(err)
	assertion.Equal(uint64(1), out.ID)
	assertion.Equal(in.Title, out.Title.String())
	assertion.Equal(in.Description, out.Description.String())
	assertion.Equal(in.Status, out.Status.String())

	repo.AssertExpectations(t)
}

func TestTaskService_GetTask_InvalidID(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	service := NewTaskService(repo)

	assertion := assert.New(t)
	_, err := service.GetTask(context.Empty(), 0)

	assertion.Error(err)
	assertion.ErrorAs(err, &ErrValidation)

	// репозиторий не должен быть вызван
	repo.AssertExpectations(t)
}

func TestTaskService_ListTasks_NormalizesLimitOffset(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	service := NewTaskService(repo)

	limit := uint64(20)
	offset := uint64(10)

	repo.
		On("List", mock.Anything, limit, offset).
		Return([]task.Task{}, nil).
		Once()

	assertion := assert.New(t)
	_, err := service.ListTasks(context.Empty(), 0, offset)
	assertion.NoError(err)

	repo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_Success(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	svc := NewTaskService(repo)

	repo.
		On("Delete", mock.Anything, uint64(1)).
		Return(nil).
		Once()

	assertion := assert.New(t)
	err := svc.DeleteTask(context.Empty(), uint64(1))

	assertion.NoError(err)

	repo.AssertExpectations(t)
}
