package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/internal/transport/http/error"
	"github.com/gin-gonic/gin"
)

func Auth(tokens service.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			error.SetError(c, http.StatusUnauthorized, errors.New("no authorization header provided"))
			c.Abort()

			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			error.SetError(c, http.StatusUnauthorized, errors.New("authorization header is invalid"))
			c.Abort()

			return
		}

		token := strings.TrimPrefix(header, prefix)
		if token == "" {
			error.SetError(c, http.StatusUnauthorized, errors.New("authorization token is empty"))
			c.Abort()

			return
		}

		userID, err := tokens.Parse(token)
		if err != nil {
			error.SetError(c, http.StatusUnauthorized, errors.New("authorization token is invalid"))
			c.Abort()

			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
