package task

import (
	"errors"
	"strconv"
	"time"

	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/Nikkoz/task-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	DefaultValueForLimit uint64 = 10
	MaxValueForLimit     uint64 = 100
)

type (
	id struct {
		Value uint64 `json:"id" uri:"id" binding:"required"`
	}

	userID struct {
		Value uint64 `json:"user_id" uri:"user_id" binding:"required"`
	}

	request struct {
		UserID      uint64     `json:"user_id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		DueDate     *time.Time `json:"due_date"` // allow null
	}
)

func getId(c *gin.Context) (*id, error) {
	i := &id{}
	if err := c.ShouldBindUri(i); err != nil {
		return nil, logger.ErrorWithContext(context.New(c), err)
	}

	return i, nil
}

func getUserID(c *gin.Context) (*userID, error) {
	value, ok := c.Get("user_id")
	if !ok {
		return nil, logger.ErrorWithContext(context.New(c), errors.New("user_id not found in context"))
	}

	uid, ok := value.(uint64)
	if !ok || uid <= 0 {
		return nil, logger.ErrorWithContext(context.New(c), errors.New("invalid user_id"))
	}

	return &userID{
		Value: uid,
	}, nil
}

func getRequest(c *gin.Context) (*request, error) {
	r := &request{}
	if err := c.ShouldBindJSON(r); err != nil {
		return nil, logger.ErrorWithContext(context.New(c), err)
	}

	return r, nil
}

func getPagination(c *gin.Context) (uint64, uint64) {
	// page number
	page, err := strconv.ParseUint(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}

	// limit of items per page
	limit, err := strconv.ParseUint(c.Query("limit"), 10, 64)
	if err != nil || limit == 0 {
		limit = DefaultValueForLimit
	}

	if limit > MaxValueForLimit {
		limit = MaxValueForLimit
	}

	return page, limit
}
