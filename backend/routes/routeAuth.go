package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/auth/register", handlers.RegisterUser)
	router.POST("/auth/login", handlers.LoginUser)
	router.GET("/auth/verifySession", handlers.VerifySessionValid)
	router.POST("/auth/logout", handlers.Logout)
	router.GET("/auth/username", handlers.FetchUsername)
	router.DELETE("/auth/deleteUser", handlers.DelAccount)
	router.POST("/auth/logout/all", handlers.LogoutAllSession)
}
