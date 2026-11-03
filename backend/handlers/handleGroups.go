package handlers

import (
	"backend/functions"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddGroup(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		functions.LogError(functions.GetFile(), sessionErr)
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "parameter name is missing"})
		return
	}
	
	createErr := functions.CreateGroup(name, userID)
	if createErr != nil {
		functions.LogError(functions.GetFile(), createErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": createErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group created Sucessfully"})
}

func FetchGroups(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		functions.LogError(functions.GetFile(), sessionErr)
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}

	groups, fetchErr := functions.GetGroups(userID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": fetchErr.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func DeleteGroup(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		functions.LogError(functions.GetFile(), sessionErr)
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}
	groupID, convErr := strconv.Atoi(c.Query("gid"))
	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(500, convErr.Error())
		return
	}

	delErr := functions.DelGroup(userID, groupID)
	if delErr != nil {
		functions.LogError(functions.GetFile(), delErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted succesfully"})

}
