package task

import "github.com/gin-gonic/gin"

func injectUserId(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}
