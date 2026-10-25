package handlers

import (
	"backend/functions"
	"log"

	"github.com/gin-gonic/gin"
)

func GetCollections(c *gin.Context) {
	collections, err := functions.FetchCollections()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, collections)
}

func NewCollection(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		log.Println("Error in handleCollections.go")
		log.Println("Wrong parameters, name is missing")
		c.JSON(400, gin.H{"message": "parameter name is missing"})
	}

	dbErr := functions.AddNewCollection(name)
	if dbErr != nil {
		log.Println("Error in handleCollections.go")
		log.Println(dbErr)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Collection created succesfully"})

}

func RemoveColellectionHandle(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		log.Println("Error in handleCollections.go")
		log.Println("Wrong parameters, name is missing")
		c.JSON(400, gin.H{"message": "parameter name is missing"})
	}

	dbErr := functions.RemoveCollection(name)
	if dbErr != nil {
		log.Println("Error in handleCollections.go")
		log.Println(dbErr)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Collection removed succesfully"})
}
