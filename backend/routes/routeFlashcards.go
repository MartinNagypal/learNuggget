package routes

import "github.com/gin-gonic/gin"
import "backend/handlers"

func FlashcardRoutes(router *gin.Engine) {
	router.POST("/addSingleCard", handlers.AddNewCard)
	router.GET("/cards", handlers.GetCards)
	router.POST("/addBulk", handlers.AddBulkCards)
	router.DELETE("/card/delete", handlers.DeleteCard)
}
