package task

import (
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

	request struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		DueDate     *time.Time `json:"due_date"` // allow null
	}

	pagination struct {
		Limit uint64 `form:"limit"`
		Page  uint64 `form:"page"`
	}
)

func getRequest(c *gin.Context) (*request, error) {
	r := request{}
	if err := c.ShouldBindJSON(&r); err != nil {
		return nil, logger.ErrorWithContext(context.New(c), err)
	}

	return &r, nil
}

func getId(c *gin.Context) (*id, error) {
	id := &id{}
	if err := c.ShouldBindUri(&id); err != nil {
		return nil, logger.ErrorWithContext(context.New(c), err)
	}

	return id, nil
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
