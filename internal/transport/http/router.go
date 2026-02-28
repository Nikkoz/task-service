package http

import (
	"net/http"

	"github.com/Nikkoz/task-service/internal/config"
	"github.com/Nikkoz/task-service/internal/transport/http/middlewares"
	"github.com/Nikkoz/task-service/internal/transport/http/task"
	"github.com/gin-gonic/gin"
)

func newRouter(taskHandler *task.Handler, isProd bool, auth config.Auth) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(middlewares.Auth(auth))
	router.Use(middlewares.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	tasks(router.Group("/tasks"), taskHandler)

	return router
}

func tasks(router *gin.RouterGroup, h *task.Handler) {
	router.POST("", h.Create)
	router.GET("", h.List)

	oneTask(router.Group("/:id"), h)
}

func oneTask(router *gin.RouterGroup, h *task.Handler) {
	router.GET("", h.Get)
	router.PUT("", h.Update)
	router.DELETE("", h.Delete)
}
