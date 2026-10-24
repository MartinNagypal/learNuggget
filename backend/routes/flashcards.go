package routes

import "github.com/gin-gonic/gin"
import "backend/handlers"

func FlashcardRoutes(router *gin.Engine) {
	router.GET("/getCollections", handlers.GetCards)
}
