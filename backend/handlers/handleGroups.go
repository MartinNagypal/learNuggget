package handlers

import (
	"backend/functions"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddGroup(c *gin.Context) {
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

	name, nameErr := functions.NormalizeName(c.Query("name"))
	if nameErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": nameErr.Error()})
		return
	}

	createErr := functions.CreateGroup(name, userID)
	if createErr != nil {
		functions.LogError(functions.GetFile(), createErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group created Sucessfully"})
}

func FetchGroups(c *gin.Context) {
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

	groups, fetchErr := functions.GetGroups(userID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func DeleteGroup(c *gin.Context) {
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

	delErr := functions.DelGroup(userID, groupID)
	if errors.Is(delErr, functions.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": delErr.Error()})
		return
	}
	if delErr != nil {
		functions.LogError(functions.GetFile(), delErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted succesfully"})
}
