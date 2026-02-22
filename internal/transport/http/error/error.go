package error

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	ID     uuid.UUID   `json:"id"`
	Error  string      `json:"message,omitempty"`
	Errors []string    `json:"errors,omitempty"`
	Info   interface{} `json:"info,omitempty"`
}

func SetError(c *gin.Context, statusCode int, errs ...error) {
	if len(errs) == 0 {
		return
	}

	response := Response{
		ID: uuid.New(),
	}

	for _, err := range errs {
		response.Errors = append(response.Errors, err.Error())
	}

	c.JSON(statusCode, response)
}
