package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	router.GET("/admin/stats", handlers.GetAdminStats)
	router.GET("/admin/users", handlers.GetAllUsers)
	router.DELETE("/admin/user/delete", handlers.AdminDeleteUser)
	router.POST("/admin/user/logout", handlers.AdminLogoutUser)
}
