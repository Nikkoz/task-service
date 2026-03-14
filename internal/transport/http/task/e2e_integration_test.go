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
	"github.com/Nikkoz/task-service/internal/domain/user"
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

func setupE2ERouter(uc *service.TaskService, uid uint64) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewHandler(uc)

	r.Use(injectUserId(uid))
	RegisterRoutes(r.Group("/tasks"), h)

	return r
}

func TestE2E_CreateTask(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		truncate(t, ctx, tx)

		u, _ := createUser(t, ctx, tx, "test@example.com")

		repo := postgres.NewTaskRepo(tx)
		svc := service.NewTaskService(repo)
		r := setupE2ERouter(svc, u.ID)

		reqBody := map[string]any{
			"user_id":     u.ID,
			"title":       "Buy milk",
			"description": "2 liters",
			"status":      "planned",
			"due_date":    nil,
		}

		bodyBytes, _ := json.Marshal(reqBody)
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

		u, _ := createUser(t, ctx, tx, "test@example.com")

		repo := postgres.NewTaskRepo(tx)
		svc := service.NewTaskService(repo)
		r := setupE2ERouter(svc, u.ID)

		title, _ := domain.NewTitle("Read book")
		description, _ := domain.NewDescription("DDD chapter")

		created, err := repo.Create(ctx, domain.Task{
			UserID:      u.ID,
			Title:       *title,
			Description: *description,
			Status:      domain.StatusInProgress,
			DueDate:     nil,
		})

		assertion := assert.New(t)
		assertion.NoError(err)

		url := "/tasks/" + strconv.FormatInt(int64(created.ID), 10)
		req := httptest.NewRequest(http.MethodGet, url, nil)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assertion.Equal(http.StatusOK, rec.Code)

		var resp map[string]any
		if err = json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}

		assertion.Equal(resp["id"], float64(created.ID))
		assertion.Equal(resp["title"], "Read book")
		assertion.Equal(resp["description"], "DDD chapter")
		assertion.Equal(resp["status"], "in_progress")
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

func createUser(t *testing.T, ctx context.Context, tx pgx.Tx, email string) (user.User, error) {
	t.Helper()

	repo := postgres.NewUserRepo(tx)

	e, _ := user.NewEmail(email)
	password, _ := user.NewPasswordHash("password")

	return repo.Create(ctx, user.User{
		Email:        *e,
		PasswordHash: *password,
	})
}
