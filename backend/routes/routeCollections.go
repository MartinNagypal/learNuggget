package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func CollectionRoutes(router *gin.Engine) {
	router.GET("/collections", handlers.GetCollections)
	router.POST("/collections/new", handlers.NewCollection)
	router.DELETE("/collections/delete", handlers.RemoveColellectionHandle)
}
