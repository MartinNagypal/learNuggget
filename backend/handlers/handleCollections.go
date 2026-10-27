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
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(500, convErr.Error())
		return
	}

	collections, err := functions.FetchCollections(sortState)
	if err != nil {
		functions.LogError(functions.GetFile(), err)
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
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Collection created succesfully"})

}

func RemoveColellectionHandle(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		log.Println("Error in handleCollections.go:")
		log.Println("Wrong parameters, name is missing")
		c.JSON(400, gin.H{"message": "parameter name is missing"})
	}

	dbErr := functions.RemoveCollection(name)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Collection removed succesfully"})
}

func GetCollectionLength(c *gin.Context) {
	collectionID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		functions.LogError("handleCollections.go", convErr)
		c.JSON(500, gin.H{"message": convErr.Error()})
		return
	}

	len, err := functions.GetCollectionLen(collectionID)

	if err != nil {
		functions.LogError("handleCollection.go", err)
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"len": len})
}
