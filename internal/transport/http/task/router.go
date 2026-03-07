package task

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	router.POST("", h.Create)
	router.GET("", h.List)

	taskRoute(router.Group("/:id"), h)
}

func taskRoute(router *gin.RouterGroup, h *Handler) {
	router.GET("", h.Get)
	router.PUT("", h.Update)
	router.DELETE("", h.Delete)
}
