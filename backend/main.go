package main

import (
	"fmt"

	"backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5174",
			"http://127.0.0.1:5174",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
	}))

	routes.FlashcardRoutes(router)
	routes.CollectionRoutes(router)

	err := router.Run(":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
}
