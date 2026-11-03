package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(router *gin.Engine) {
	router.POST("/group/create", handlers.AddGroup)
	router.GET("/group/fetchall", handlers.FetchGroups)
	router.DELETE("/group/delete", handlers.DeleteGroup)
}
