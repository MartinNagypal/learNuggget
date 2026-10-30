package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/auth/register", handlers.RegisterUser)
	router.POST("/auth/login", handlers.LoginUser)
	router.GET("/auth/verifySession", handlers.VerifySessionValid)
	router.GET("/auth/logout", handlers.Logout)
}
