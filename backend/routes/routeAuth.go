package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/auth/register", handlers.RegisterUser)
}
