package handlers

import (
	"backend/functions"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCollections(c *gin.Context) {
	sortState, convErr := strconv.Atoi(c.Query("orderState"))
	if convErr != nil {
		log.Println("Error in handleCollections.go:")
		log.Println(convErr)
		c.JSON(500, convErr.Error())
		return
	}

	collections, err := functions.FetchCollections(sortState)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
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
