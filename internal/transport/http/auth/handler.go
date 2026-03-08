package auth

import (
	"errors"
	"net/http"

	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/internal/service"
	httpError "github.com/Nikkoz/task-service/internal/transport/http/error"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Register(c *gin.Context) {
	req, err := getRequest(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	var ctx = context.New(c)

	in := service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := h.service.Register(ctx, in)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, repository.ErrAlreadyExists) {
			code = http.StatusConflict
		}

		httpError.SetError(c, code, err)

		return
	}

	c.JSON(http.StatusCreated, toRegisterResponse(out))
}

func (h *Handler) Login(c *gin.Context) {
	req, err := getRequest(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	var ctx = context.New(c)

	in := service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := h.service.Login(ctx, in)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidCredentials) {
			code = http.StatusUnauthorized
		} else if errors.Is(err, service.ErrValidation) {
			code = http.StatusBadRequest
		}

		httpError.SetError(c, code, err)

		return
	}

	c.JSON(http.StatusOK, toLoginResponse(out))
}
