package handlers

import (
	"backend/functions"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// requireCollectionOwner checks that the collection belongs to the user. If it
// does not (or the check fails), it writes the error response and returns false.
func requireCollectionOwner(c *gin.Context, collectionID, userID int) bool {
	owns, ownErr := functions.UserOwnsCollection(collectionID, userID)
	if ownErr != nil {
		functions.LogError(functions.GetFile(), ownErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return false
	}
	if !owns {
		c.JSON(http.StatusNotFound, gin.H{"message": functions.ErrCollectionNotFound.Error()})
		return false
	}
	return true
}

func GetCollections(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	groupID, convErr := strconv.Atoi(c.Query("gid"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No group id given"})
		return
	}

	collections, fetchErr := functions.FetchCollections(userID, groupID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, collections)
}

func NewCollection(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	groupID, convErr := strconv.Atoi(c.Query("gid"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No group id given"})
		return
	}
	name, nameErr := functions.NormalizeName(c.Query("name"))
	if nameErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": nameErr.Error()})
		return
	}

	addErr := functions.AddNewCollection(name, userID, groupID)
	if errors.Is(addErr, functions.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": addErr.Error()})
		return
	}
	if addErr != nil {
		functions.LogError(functions.GetFile(), addErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection created succesfully"})
}

func RemoveColellectionHandle(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	collectionID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No collection id given"})
		return
	}

	removeErr := functions.RemoveCollection(collectionID, userID)
	if errors.Is(removeErr, functions.ErrCollectionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": removeErr.Error()})
		return
	}
	if removeErr != nil {
		functions.LogError(functions.GetFile(), removeErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Collection removed succesfully"})
}

func GetCollectionLength(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	collectionID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No collection id given"})
		return
	}
	if !requireCollectionOwner(c, collectionID, userID) {
		return
	}

	length, lenErr := functions.GetCollectionLen(collectionID)
	if lenErr != nil {
		functions.LogError(functions.GetFile(), lenErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"len": length})
}

func GetCollectionStats(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	collectionID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No collection id given"})
		return
	}
	if !requireCollectionOwner(c, collectionID, userID) {
		return
	}
	stats, statsErr := functions.GetCollectionStats(collectionID)
	if statsErr != nil {
		functions.LogError(functions.GetFile(), statsErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func UpdateCollectionProgress(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}
	var result functions.LearnResult
	bindErr := c.ShouldBindJSON(&result)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	if !requireCollectionOwner(c, result.CollectionID, userID) {
		return
	}
	response, updateErr := functions.UpdateCollectionStats(result)
	if errors.Is(updateErr, functions.ErrInvalidLearnResult) {
		c.JSON(http.StatusBadRequest, gin.H{"message": updateErr.Error()})
		return
	}
	if updateErr != nil {
		functions.LogError(functions.GetFile(), updateErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetMarkedCollections(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	userID, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	collections, fetchErr := functions.FetchMarkedCollections(userID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, collections)
}
