//go:build integration

package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	domain "github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/repository/postgres"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/internal/testutil"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("ENV_FILE", ".testing")

	code := m.Run()

	testutil.ClosePool()
	os.Exit(code)
}

func setupE2ERouter(uc *service.TaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewHandler(uc)

	RegisterRoutes(r.Group("/tasks"), h)

	return r
}

func TestE2E_CreateTask(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		truncate(t, ctx, tx)

		repo := postgres.NewTaskRepo(tx)
		svc := service.NewTaskService(repo)
		r := setupE2ERouter(svc)

		reqBody := map[string]any{
			"title":       "Buy milk",
			"description": "2 liters",
			"status":      "planned",
			"due_date":    nil,
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assertion := assert.New(t)
		assertion.Equal(http.StatusCreated, rec.Code)

		var resp map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}

		assertion.NotEmpty(resp["id"])
		assertion.Equal("Buy milk", resp["title"])
		assertion.Equal("2 liters", resp["description"])
		assertion.Equal("planned", resp["status"])
		assertion.Equal(nil, resp["due_date"])
	})
}

func TestE2E_GetTask(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		truncate(t, ctx, tx)

		repo := postgres.NewTaskRepo(tx)
		svc := service.NewTaskService(repo)
		r := setupE2ERouter(svc)

		// создаём запись напрямую через repo, чтобы потом получить её по GET
		title, err := domain.NewTitle("Read book")
		if err != nil {
			t.Fatalf("new title: %v", err)
		}

		description, err := domain.NewDescription("DDD chapter")
		if err != nil {
			t.Fatalf("new description: %v", err)
		}

		created, err := repo.Create(ctx, domain.Task{
			Title:       *title,
			Description: *description,
			Status:      domain.StatusInProgress,
			DueDate:     nil,
		})
		if err != nil {
			t.Fatalf("create task: %v", err)
		}

		url := "/tasks/" + strconv.FormatInt(int64(created.ID), 10)
		req := httptest.NewRequest(http.MethodGet, url, nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var resp map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}

		if resp["id"] != float64(created.ID) {
			t.Fatalf("expected id %d, got %v", created.ID, resp["id"])
		}

		if resp["title"] != "Read book" {
			t.Fatalf("expected title %q, got %v", "Read book", resp["title"])
		}

		if resp["status"] != "in_progress" {
			t.Fatalf("expected status %q, got %v", "in_progress", resp["status"])
		}
	})
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
