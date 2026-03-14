package task

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

func (h *Handler) Create(c *gin.Context) {
	req, err := getRequest(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	var ctx = context.New(c)

	out, err := h.service.CreateTask(ctx, service.CreateTaskInput{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
	})
	if err != nil {
		httpError.SetError(c, http.StatusInternalServerError, err)

		return
	}

	//if out == nil {
	//	return
	//}

	c.JSON(http.StatusCreated, toResponse(out))
}

func (h *Handler) Update(c *gin.Context) {
	taskId, err := getId(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	userId, err := getUserID(c)
	if err != nil {
		httpError.SetError(c, http.StatusUnauthorized, err)

		return
	}

	req, err := getRequest(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	var ctx = context.New(c)

	out, err := h.service.UpdateTask(ctx, taskId.Value, userId.Value, service.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
	})
	if err != nil {
		httpError.SetError(c, http.StatusInternalServerError, err)

		return
	}

	//if out == nil {
	//	return
	//}

	c.JSON(http.StatusOK, toResponse(out))
}

func (h *Handler) Get(c *gin.Context) {
	taskId, err := getId(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	userId, err := getUserID(c)
	if err != nil {
		httpError.SetError(c, http.StatusUnauthorized, err)

		return
	}

	var ctx = context.New(c)

	out, err := h.service.GetTask(ctx, taskId.Value, userId.Value)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			code = http.StatusNotFound
		}

		httpError.SetError(c, code, err)

		return
	}

	c.JSON(http.StatusOK, toResponse(out))
}

func (h *Handler) List(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		httpError.SetError(c, http.StatusUnauthorized, err)

		return
	}

	page, limit := getPagination(c)
	if page > 0 {
		page = page - 1
	}

	var ctx = context.New(c)

	out, err := h.service.ListTasks(ctx, userId.Value, limit, page*limit)
	if err != nil {
		httpError.SetError(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, toListResponse(limit, page, out))
}

func (h *Handler) Delete(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		httpError.SetError(c, http.StatusUnauthorized, err)

		return
	}

	taskId, err := getId(c)
	if err != nil {
		httpError.SetError(c, http.StatusBadRequest, err)

		return
	}

	var ctx = context.New(c)

	if err = h.service.DeleteTask(ctx, taskId.Value, userId.Value); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			code = http.StatusNotFound
		}
		httpError.SetError(c, code, err)

		return
	}

	c.Status(http.StatusNoContent)
}
