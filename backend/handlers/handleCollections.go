package handlers

import (
	"backend/functions"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCollections(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}

	sortState, convErr := strconv.Atoi(c.Query("orderState"))
	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(500, convErr.Error())
		return
	}

	collections, err := functions.FetchCollections(sortState, userID)
	if err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, collections)
}

func NewCollection(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}

	name := c.Query("name")
	if name == "" {
		log.Println("Error in handleCollections.go")
		log.Println("Wrong parameters, name is missing")
		c.JSON(400, gin.H{"message": "parameter name is missing"})
	}

	dbErr := functions.AddNewCollection(name, userID)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Collection created succesfully"})

}

func RemoveColellectionHandle(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}

	name := c.Query("name")
	if name == "" {
		log.Println("Error in handleCollections.go:")
		log.Println("Wrong parameters, name is missing")
		c.JSON(400, gin.H{"message": "parameter name is missing"})
	}

	dbErr := functions.RemoveCollection(name, userID)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Collection removed succesfully"})
}

func GetCollectionLength(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	_, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}

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
