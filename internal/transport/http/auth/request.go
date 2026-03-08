package auth

import (
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/Nikkoz/task-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func getRequest(c *gin.Context) (*request, error) {
	r := request{}
	if err := c.ShouldBindJSON(&r); err != nil {
		return nil, logger.ErrorWithContext(context.New(c), err)
	}

	return &r, nil
}
