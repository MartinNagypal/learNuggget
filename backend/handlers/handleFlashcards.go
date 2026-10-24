package handlers

import "github.com/gin-gonic/gin"

func GetCards(c *gin.Context) {
	c.JSON(200, gin.H{"car1": "hallo"})
}
