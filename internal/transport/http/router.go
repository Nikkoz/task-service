package http

import (
	"net/http"

	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/internal/transport/http/auth"
	"github.com/Nikkoz/task-service/internal/transport/http/middlewares"
	"github.com/Nikkoz/task-service/internal/transport/http/task"
	"github.com/gin-gonic/gin"
)

func newRouter(taskHandler *task.Handler, authHandler *auth.Handler, tokenManager service.TokenManager, isProd bool) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	auth.RegisterRoutes(router.Group("/auth"), authHandler)

	protected := router.Group("")
	protected.Use(middlewares.Auth(tokenManager))
	task.RegisterRoutes(protected.Group("/tasks"), taskHandler)

	return router
}
