//go:build integration

package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nikkoz/task-service/internal/testutil"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/repository"
)

func TestTaskRepo_GetByID_OtherUserNotFound(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		u, err := createUser(t, ctx, tx, "testA@example.com")
		assertion.NoError(err)

		created, err := createTaskForUser(t, ctx, repo, u.ID, "task-a")
		assertion.NoError(err)

		_, err = repo.GetByID(ctx, 2, created.ID)

		assertion.Error(err)
		assertion.ErrorAs(err, &repository.ErrNotFound)
	})
}

func TestTaskRepo_Update_OtherUserNotFound(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		ua, err := createUser(t, ctx, tx, "testA@example.com")
		assertion.NoError(err)

		ub, err := createUser(t, ctx, tx, "testB@example.com")
		assertion.NoError(err)

		created, err := createTaskForUser(t, ctx, repo, ua.ID, "task-a")
		assertion.NoError(err)

		newTitle, _ := task.NewTitle("updated-by-other-user")
		newDescription, _ := task.NewDescription("updated description")
		due, _ := task.NewDueDate(time.Now().Add(24 * time.Hour))

		_, err = repo.Update(ctx, task.Task{
			ID:          created.ID,
			UserID:      ub.ID, // another user
			Title:       *newTitle,
			Description: *newDescription,
			Status:      task.StatusDone,
			DueDate:     due,
		})

		assertion.Error(err)
		assertion.ErrorAs(err, &repository.ErrNotFound)

		found, err := repo.GetByID(ctx, created.ID, ua.ID)

		assertion.NoError(err)
		assertion.Equal(created.Title.String(), found.Title.String())
	})
}

func TestTaskRepo_Delete_OtherUserNotFound(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		ua, err := createUser(t, ctx, tx, "testA@example.com")
		assertion.NoError(err)

		ub, err := createUser(t, ctx, tx, "testB@example.com")
		assertion.NoError(err)

		created, err := createTaskForUser(t, ctx, repo, ua.ID, "task-a")
		assertion.NoError(err)

		err = repo.Delete(ctx, created.ID, ub.ID) // another user

		assertion.Error(err)
		assertion.ErrorAs(err, &repository.ErrNotFound)

		// владелец всё ещё видит задачу
		_, err = repo.GetByID(ctx, created.ID, ua.ID)
		assertion.NoError(err)
	})
}

func TestTaskRepo_List_OnlyOwnerTasks(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		ua, err := createUser(t, ctx, tx, "testA@example.com")
		assertion.NoError(err)

		ub, err := createUser(t, ctx, tx, "testB@example.com")
		assertion.NoError(err)

		a1, err := createTaskForUser(t, ctx, repo, ua.ID, "user1-task-1")
		a2, err := createTaskForUser(t, ctx, repo, ua.ID, "user1-task-2")
		_, err = createTaskForUser(t, ctx, repo, ub.ID, "user2-task-1")

		got, err := repo.List(ctx, ua.ID, 10, 0)
		assertion.NoError(err)
		assertion.NotEmpty(got)
		assertion.Len(got, 2)

		expectedIDs := map[uint64]struct{}{
			a1.ID: {},
			a2.ID: {},
		}

		for _, item := range got {
			assertion.Equal(ua.ID, item.UserID)
			assertion.Contains(expectedIDs, item.ID)
		}
	})
}

func createTaskForUser(t *testing.T, ctx context.Context, repo *TaskRepo, userID uint64, titleRaw string) (task.Task, error) {
	t.Helper()

	title, _ := task.NewTitle(titleRaw)
	description, _ := task.NewDescription(fmt.Sprintf("desc-%s", titleRaw))
	due, _ := task.NewDueDate(time.Now().Add(24 * time.Hour))

	return repo.Create(ctx, task.Task{
		UserID:      userID,
		Title:       *title,
		Description: *description,
		Status:      task.StatusPlanned,
		DueDate:     due,
	})
}
