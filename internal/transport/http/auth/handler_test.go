package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/internal/transport/http/auth/mocks"
)

func setupRouter(svc *mocks.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewHandler(svc)

	RegisterRoutes(r.Group("/auth"), h)

	return r
}

func TestRegister_Success(t *testing.T) {
	svc := mocks.NewService(t)
	router := setupRouter(svc)

	email := "test@example.com"
	password := "secret123"

	login, _ := user.NewEmail(email)
	svc.
		On(
			"Register",
			mock.Anything,
			mock.MatchedBy(func(in service.RegisterInput) bool {
				return in.Email == login.String()
			}),
		).
		Return(user.User{
			ID:    1,
			Email: *login,
		}, nil).
		Once()

	body := request{
		Email:    email,
		Password: password,
	}

	data, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusCreated, rec.Code)
}

func TestRegister_Duplicate(t *testing.T) {
	svc := mocks.NewService(t)
	router := setupRouter(svc)

	svc.
		On("Register", mock.Anything, mock.Anything).
		Return(user.User{}, repository.ErrAlreadyExists).
		Once()

	email := "test@example.com"
	password := "secret123"
	body := request{
		Email:    email,
		Password: password,
	}

	data, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusConflict, rec.Code)
}

func TestRegister_InvalidJSON(t *testing.T) {
	svc := mocks.NewService(t)
	router := setupRouter(svc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewBuffer([]byte("{invalid json")),
	)

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}

func TestLogin_Success(t *testing.T) {
	svc := mocks.NewService(t)
	router := setupRouter(svc)

	email := "test@example.com"
	password := "secret123"
	svc.
		On("Login", mock.Anything, service.LoginInput{
			Email:    email,
			Password: password,
		}).
		Return(service.LoginResult{
			Token: "jwt-token",
		}, nil).
		Once()

	body := request{
		Email:    email,
		Password: password,
	}

	data, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusOK, rec.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	svc := mocks.NewService(t)
	router := setupRouter(svc)

	svc.
		On("Login", mock.Anything, mock.Anything).
		Return(service.LoginResult{}, service.ErrInvalidCredentials).
		Once()

	body := request{
		Email:    "test@example.com",
		Password: "wrong",
	}

	data, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusUnauthorized, rec.Code)
}

func TestLogin_ErrValidation(t *testing.T) {
	svc := mocks.NewService(t)
	router := setupRouter(svc)

	svc.
		On("Login", mock.Anything, mock.Anything).
		Return(service.LoginResult{}, service.ErrValidation).
		Once()

	body := request{
		Email:    "test@example.com",
		Password: "wrong",
	}

	data, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assertion := assert.New(t)
	assertion.Equal(http.StatusBadRequest, rec.Code)
}
