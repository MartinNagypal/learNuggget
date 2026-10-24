package handlers

import "github.com/gin-gonic/gin"

func GetCollections(c *gin.Context) {
	collections := []string{"Algorithmen", "Mathe", "Komplexität"}
	c.JSON(200, collections)
}
