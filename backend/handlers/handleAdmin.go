package handlers

import (
	"backend/functions"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAdminStats(c *gin.Context) {
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
	root, rootErr := functions.IsRoot(userID)
	if rootErr != nil {
		functions.LogError(functions.GetFile(), rootErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if !root {
		c.JSON(http.StatusForbidden, gin.H{"message": "No permission"})
		return
	}

	stats, statsErr := functions.GetAdminStats()
	if statsErr != nil {
		functions.LogError(functions.GetFile(), statsErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GetAllUsers(c *gin.Context) {
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
	root, rootErr := functions.IsRoot(userID)
	if rootErr != nil {
		functions.LogError(functions.GetFile(), rootErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if !root {
		c.JSON(http.StatusForbidden, gin.H{"message": "No permission"})
		return
	}

	users, usersErr := functions.FetchAllUsers()
	if usersErr != nil {
		functions.LogError(functions.GetFile(), usersErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func AdminDeleteUser(c *gin.Context) {
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
	root, rootErr := functions.IsRoot(userID)
	if rootErr != nil {
		functions.LogError(functions.GetFile(), rootErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if !root {
		c.JSON(http.StatusForbidden, gin.H{"message": "No permission"})
		return
	}

	targetID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user id given"})
		return
	}

	delErr := functions.DeleteAccount(targetID)
	if errors.Is(delErr, functions.ErrRootCannotBeDeleted) {
		c.JSON(http.StatusForbidden, gin.H{"message": delErr.Error()})
		return
	}
	if errors.Is(delErr, functions.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": delErr.Error()})
		return
	}
	if delErr != nil {
		functions.LogError(functions.GetFile(), delErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func AdminLogoutUser(c *gin.Context) {
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
	root, rootErr := functions.IsRoot(userID)
	if rootErr != nil {
		functions.LogError(functions.GetFile(), rootErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if !root {
		c.JSON(http.StatusForbidden, gin.H{"message": "No permission"})
		return
	}

	targetID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user id given"})
		return
	}

	exists, existsErr := functions.UserExists(targetID)
	if existsErr != nil {
		functions.LogError(functions.GetFile(), existsErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": functions.ErrUserNotFound.Error()})
		return
	}

	targetRoot, targetRootErr := functions.IsRoot(targetID)
	if targetRootErr != nil {
		functions.LogError(functions.GetFile(), targetRootErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if targetRoot {
		c.JSON(http.StatusForbidden, gin.H{"message": "Root user can't be logged out here"})
		return
	}

	logoutErr := functions.LogoutAllSessions(targetID)
	if logoutErr != nil {
		functions.LogError(functions.GetFile(), logoutErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged out"})
}
