package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func SocialRoutes(router *gin.Engine) {
	router.GET("/user/score", handlers.GetUserScore)
	router.GET("/user/top", handlers.GetTopUsers)
}
