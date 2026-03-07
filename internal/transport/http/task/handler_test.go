package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Nikkoz/task-service/internal/domain/task"
	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/internal/transport/http/task/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(uc *mocks.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewHandler(uc)

	RegisterRoutes(r.Group("/tasks"), h)

	return r
}

func TestCreate_InvalidJSON(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title":`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}

func TestCreate_Error(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	titleRaw := "Buy milk"
	descriptionRaw := "2 liters"
	statusRaw := "planned"
	dueDateRaw := time.Now().UTC().Add(24 * time.Hour)

	svc.
		On(
			"CreateTask",
			mock.Anything,
			mock.MatchedBy(func(in service.CreateTaskInput) bool {
				return in.Title == titleRaw &&
					in.Description == descriptionRaw &&
					in.Status == statusRaw &&
					in.DueDate.Format(time.RFC3339) == dueDateRaw.Format(time.RFC3339)
			}),
		).
		Return(task.Task{}, errors.New("create error")).
		Once()

	reqBody := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		DueDate     string `json:"due_date"`
	}{
		Title:       titleRaw,
		Description: descriptionRaw,
		Status:      statusRaw,
		DueDate:     dueDateRaw.Format(time.RFC3339),
	}

	assertion := assert.New(t)
	bodyBytes, err := json.Marshal(reqBody)

	assertion.NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion.Equal(http.StatusInternalServerError, rec.Code)

	svc.AssertExpectations(t)
}

func TestCreate_Success(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	titleRaw := "Buy milk"
	descriptionRaw := "2 liters"
	statusRaw := "planned"
	dueDateRaw := time.Now().UTC().Add(24 * time.Hour)

	title, _ := task.NewTitle(titleRaw)
	description, _ := task.NewDescription(descriptionRaw)
	status, _ := task.NewStatus(statusRaw)
	dueDate, _ := task.NewDueDate(dueDateRaw)

	created := task.Task{
		ID:          1,
		Title:       *title,
		Description: *description,
		Status:      *status,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	svc.
		On(
			"CreateTask",
			mock.Anything,
			mock.MatchedBy(func(in service.CreateTaskInput) bool {
				return in.Title == titleRaw &&
					in.Description == descriptionRaw &&
					in.Status == statusRaw &&
					in.DueDate.Format(time.RFC3339) == dueDateRaw.Format(time.RFC3339)
			}),
		).
		Return(created, nil).
		Once()

	reqBody := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		DueDate     string `json:"due_date"`
	}{
		Title:       titleRaw,
		Description: descriptionRaw,
		Status:      statusRaw,
		DueDate:     dueDateRaw.Format(time.RFC3339),
	}

	assertion := assert.New(t)
	bodyBytes, err := json.Marshal(reqBody)

	assertion.NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion.Equal(http.StatusCreated, rec.Code)

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	assertion.Equal(created.ID, uint64(resp["id"].(float64)))
	assertion.Equal(titleRaw, resp["title"])
	assertion.Equal(descriptionRaw, resp["description"])
	assertion.Equal(statusRaw, resp["status"])
	assertion.Equal(dueDateRaw.String(), resp["due_date"])

	svc.AssertExpectations(t)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/1",
		strings.NewReader(`{"title":`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}

func TestUpdate_InvalidId(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/one",
		strings.NewReader(`{"title":"test"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}

func TestUpdate_EmptyId(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks",
		strings.NewReader(`{"title":"test"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusNotFound, rec.Code)
}

func TestUpdate_Error(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)
	titleRaw := "Buy milk"
	descriptionRaw := "2 liters"
	statusRaw := "planned"
	dueDateRaw := time.Now().UTC().Add(24 * time.Hour)

	svc.
		On(
			"UpdateTask",
			mock.Anything,
			taskID,
			mock.MatchedBy(func(in service.UpdateTaskInput) bool {
				return in.Title == titleRaw &&
					in.Description == descriptionRaw &&
					in.Status == statusRaw &&
					in.DueDate.Format(time.RFC3339) == dueDateRaw.Format(time.RFC3339)
			}),
		).
		Return(task.Task{}, errors.New("update error")).
		Once()

	reqBody := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		DueDate     string `json:"due_date"`
	}{
		Title:       titleRaw,
		Description: descriptionRaw,
		Status:      statusRaw,
		DueDate:     dueDateRaw.Format(time.RFC3339),
	}

	assertion := assert.New(t)
	bodyBytes, err := json.Marshal(reqBody)

	assertion.NoError(err)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion.Equal(http.StatusInternalServerError, rec.Code)

	svc.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)
	titleRaw := "Buy milk"
	descriptionRaw := "2 liters"
	statusRaw := "planned"
	dueDateRaw := time.Now().UTC().Add(24 * time.Hour)

	title, _ := task.NewTitle(titleRaw)
	description, _ := task.NewDescription(descriptionRaw)
	status, _ := task.NewStatus(statusRaw)
	dueDate, _ := task.NewDueDate(dueDateRaw)

	updated := task.Task{
		ID:          taskID,
		Title:       *title,
		Description: *description,
		Status:      *status,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	svc.
		On(
			"UpdateTask",
			mock.Anything,
			taskID,
			mock.MatchedBy(func(in service.UpdateTaskInput) bool {
				return in.Title == titleRaw &&
					in.Description == descriptionRaw &&
					in.Status == statusRaw &&
					in.DueDate.Format(time.RFC3339) == dueDateRaw.Format(time.RFC3339)
			}),
		).
		Return(updated, nil).
		Once()

	reqBody := struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		DueDate     string `json:"due_date"`
	}{
		Title:       titleRaw,
		Description: descriptionRaw,
		Status:      statusRaw,
		DueDate:     dueDateRaw.Format(time.RFC3339),
	}

	assertion := assert.New(t)
	bodyBytes, err := json.Marshal(reqBody)

	assertion.NoError(err)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion.Equal(http.StatusOK, rec.Code)

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	assertion.Equal(updated.ID, uint64(resp["id"].(float64)))
	assertion.Equal(titleRaw, resp["title"])
	assertion.Equal(descriptionRaw, resp["description"])
	assertion.Equal(statusRaw, resp["status"])
	assertion.Equal(dueDateRaw.String(), resp["due_date"])

	svc.AssertExpectations(t)
}

func TestGet_InvalidId(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/one",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}

func TestGet_NotFound(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)

	svc.
		On("GetTask", mock.Anything, taskID).
		Return(task.Task{}, repository.ErrNotFound).
		Once()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusNotFound, rec.Code)
}

func TestGet_Error(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)

	svc.
		On("GetTask", mock.Anything, taskID).
		Return(task.Task{}, errors.New("get error")).
		Once()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusInternalServerError, rec.Code)
}

func TestGet_Success(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)
	titleRaw := "Buy milk"
	descriptionRaw := "2 liters"
	statusRaw := "planned"
	dueDateRaw := time.Now().UTC().Add(24 * time.Hour)

	title, _ := task.NewTitle(titleRaw)
	description, _ := task.NewDescription(descriptionRaw)
	status, _ := task.NewStatus(statusRaw)
	dueDate, _ := task.NewDueDate(dueDateRaw)

	got := task.Task{
		ID:          taskID,
		Title:       *title,
		Description: *description,
		Status:      *status,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	svc.
		On("GetTask", mock.Anything, taskID).
		Return(got, nil).
		Once()

	assertion := assert.New(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion.Equal(http.StatusOK, rec.Code)

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	assertion.Equal(got.ID, uint64(resp["id"].(float64)))
	assertion.Equal(titleRaw, resp["title"])
	assertion.Equal(descriptionRaw, resp["description"])
	assertion.Equal(statusRaw, resp["status"])
	assertion.Equal(dueDateRaw.String(), resp["due_date"])

	svc.AssertExpectations(t)
}

func TestDelete_InvalidId(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/one",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}

func TestDelete_NotFound(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)

	svc.
		On("DeleteTask", mock.Anything, taskID).
		Return(repository.ErrNotFound).
		Once()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusNotFound, rec.Code)
}

func TestDelete_Error(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)

	svc.
		On("DeleteTask", mock.Anything, taskID).
		Return(errors.New("delete error")).
		Once()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusInternalServerError, rec.Code)
}

func TestDelete_Success(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	taskID := uint64(1)

	svc.
		On("DeleteTask", mock.Anything, taskID).
		Return(nil).
		Once()

	assertion := assert.New(t)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/"+strconv.FormatUint(taskID, 10),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion.Equal(http.StatusNoContent, rec.Code)

	svc.AssertExpectations(t)
}

func TestList_Error(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	svc.
		On("ListTasks", mock.Anything, uint64(10), uint64(10)).
		Return([]task.Task{}, errors.New("list error")).
		Once()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks?page=2&limit=10",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusInternalServerError, rec.Code)
}

func TestList_Success(t *testing.T) {
	svc := mocks.NewService(t)
	r := setupRouter(svc)

	svc.
		On("ListTasks", mock.Anything, uint64(10), uint64(0)).
		Return([]task.Task{}, nil).
		Once()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks?page=0&limit=0",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusOK, rec.Code)
}
