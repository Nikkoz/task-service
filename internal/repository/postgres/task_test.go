//go:build integration

package postgres

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/testutil"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("ENV_FILE", ".testing")

	code := m.Run()

	testutil.ClosePool()
	os.Exit(code)
}

func TestTaskRepo_CreateGet(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		title, _ := task.NewTitle("test")
		descr, _ := task.NewDescription("test description")
		due, _ := task.NewDueDate(time.Now().Add(24 * time.Hour))
		model := task.Task{
			Title:       *title,
			Description: *descr,
			Status:      task.StatusPlanned,
			DueDate:     due,
		}

		created, err := repo.Create(ctx, model)
		assertion.NoError(err)

		got, err := repo.GetByID(ctx, created.ID)

		assertion.NoError(err)
		assertion.Equal(model.Title, got.Title)
		assertion.Equal(model.Description, got.Description)
		assertion.Equal(model.Status, got.Status)
		assertion.Equal(
			model.DueDate.DateTime().Format("2020-01-01 10:00:00"),
			got.DueDate.DateTime().Format("2020-01-01 10:00:00"),
		)
	})
}

func TestTaskRepo_List_LimitOffset(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		truncate(t, ctx, tx)

		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		created := make([]task.Task, 0, 3)

		for i := 1; i <= 3; i++ {
			title, _ := task.NewTitle(fmt.Sprintf("list-%d", i))
			desc, _ := task.NewDescription("description")
			due, _ := task.NewDueDate(time.Now().Add(time.Duration(i) * time.Hour))

			out, err := repo.Create(ctx, task.Task{
				Title:       *title,
				Description: *desc,
				Status:      task.StatusPlanned,
				DueDate:     due,
			})
			assertion.NoError(err)

			created = append(created, out)
		}

		expectedAll := reverseIDs(created)
		cases := []struct {
			name          string
			limit, offset uint64
			wantIdx       []uint64 // индексы из created
		}{
			{"all", 10, 0, expectedAll},
			{"page1", 2, 0, expectedAll[:2]},
			{"page2", 2, 2, expectedAll[2:]},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				got, err := repo.List(ctx, c.limit, c.offset)
				assertion.NoError(err)
				assertion.Len(got, len(c.wantIdx))
				assertion.Equal(c.wantIdx, taskIDs(got))
			})
		}
	})
}

func TestTaskRepo_Update(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		title, _ := task.NewTitle("before")
		desc, _ := task.NewDescription("before-desc")
		due, _ := task.NewDueDate(time.Now().Add(time.Hour))

		created, err := repo.Create(ctx, task.Task{
			Title:       *title,
			Description: *desc,
			Status:      task.StatusPlanned,
			DueDate:     due,
		})
		assertion.NoError(err)

		newTitle, _ := task.NewTitle("after")
		newDesc, _ := task.NewDescription("after-desc")
		newDue, _ := task.NewDueDate(time.Now().UTC().Add(48 * time.Hour))

		updated, err := repo.Update(ctx, task.Task{
			ID:          created.ID,
			Title:       *newTitle,
			Description: *newDesc,
			Status:      task.StatusDone,
			DueDate:     newDue,
		})
		assertion.NoError(err)
		assertion.NotEqual(created.Title, updated.Title)
		assertion.NotEqual(created.Description, updated.Description)
		assertion.NotEqual(created.Status, updated.Status)
		assertion.NotEqual(
			created.DueDate.DateTime().Format("2020-01-01 10:00:00"),
			updated.DueDate.DateTime().Format("2020-01-01 10:00:00"),
		)

		// verify persisted
		got, err := repo.GetByID(ctx, created.ID)
		assertion.NoError(err)
		assertion.Equal(updated.ID, got.ID)
		assertion.Equal(updated.Title, got.Title)
		assertion.Equal(updated.Description, got.Description)
		assertion.Equal(updated.Status, got.Status)
		assertion.Equal(
			got.DueDate.DateTime().Format("2020-01-01 10:00:00"),
			got.DueDate.DateTime().Format("2020-01-01 10:00:00"),
		)
	})
}

func TestTaskRepo_Delete(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewTaskRepo(tx)
		assertion := assert.New(t)

		title, _ := task.NewTitle("to-delete")
		desc, _ := task.NewDescription("d")

		created, err := repo.Create(ctx, task.Task{
			Title:       *title,
			Description: *desc,
			Status:      task.StatusPlanned,
			DueDate:     nil,
		})
		assertion.NoError(err)

		err = repo.Delete(ctx, created.ID)
		assertion.NoError(err)

		_, err = repo.GetByID(ctx, created.ID)
		assertion.Error(err)
	})
}

func reverseIDs(tasks []task.Task) []uint64 {
	ids := make([]uint64, len(tasks))
	for i := range tasks {
		ids[i] = tasks[len(tasks)-1-i].ID
	}

	return ids
}

func taskIDs(ts []task.Task) []uint64 {
	ids := make([]uint64, len(ts))
	for i := range ts {
		ids[i] = ts[i].ID
	}
	return ids
}

func truncate(t *testing.T, ctx context.Context, tx pgx.Tx) {
	_, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(424242)`)
	if err != nil {
		t.Fatalf("pg_advisory_xact_lock: %v", err)
	}

	_, err = tx.Exec(ctx, `TRUNCATE TABLE tasks RESTART IDENTITY`)
	if err != nil {
		t.Fatalf("truncate tasks: %v", err)
	}
}
