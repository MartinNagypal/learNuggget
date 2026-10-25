package handlers

import (
	"backend/functions"
	"log"

	"github.com/gin-gonic/gin"
)

func GetCollections(c *gin.Context) {
	collections, err := functions.FetchCollections()
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, collections)
}
