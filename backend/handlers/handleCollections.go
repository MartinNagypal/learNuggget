package handlers

import (
	"backend/functions"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

	groupID, convErr := strconv.Atoi(c.Query("gid"))
	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(500, gin.H{"message": convErr.Error()})
		return
	}

	collections, err := functions.FetchCollections(userID, groupID)
	if err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
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

	name := strings.TrimSpace(c.Query("name"))
	groupID, convErr := strconv.Atoi(c.Query("gid"))
	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(500, gin.H{"message": convErr.Error()})
		return
	}
	if name == "" {
		functions.LogError(functions.GetFile(), errors.New("Name is empty"))
		c.JSON(400, gin.H{"message": "parameter name is missing"})
		return
	}

	dbErr := functions.AddNewCollection(name, userID, groupID)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr})
		return
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

	collectionID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(500, convErr.Error())
		return
	}

	dbErr := functions.RemoveCollection(collectionID, userID)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr})
		return
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

func GetCollectionStats(c *gin.Context) {
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
	stats, statsErr := functions.GetCollectionStats(collectionID)
	if statsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": statsErr.Error()})
		return
	}
	c.JSON(200, stats)
}

func UpdateCollectionProgress(c *gin.Context) {
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
	var updatedStats functions.UpdatedCollectionStatsElement
	bindErr := c.ShouldBind(&updatedStats)
	if bindErr != nil {
		functions.LogError(functions.GetFile(), bindErr)
		c.JSON(500, gin.H{"message": bindErr.Error()})
		return
	}
	updateErr := functions.UpdateCollectionStats(updatedStats)
	if updateErr != nil {
		functions.LogError(functions.GetFile(), updateErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": updateErr.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Collection updated succesfully"})
}
