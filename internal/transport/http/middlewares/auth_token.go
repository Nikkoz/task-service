package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Nikkoz/task-service/internal/config"
	"github.com/Nikkoz/task-service/internal/transport/http/error"
	"github.com/gin-gonic/gin"
)

func Auth(config config.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			error.SetError(c, http.StatusUnauthorized, errors.New("no Authorization header provided"))
			c.Abort()

			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token != config.Token {
			error.SetError(c, http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()

			return
		}

		c.Next()
	}
}
