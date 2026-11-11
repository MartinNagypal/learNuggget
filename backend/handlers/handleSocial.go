package handlers

import (
	"backend/functions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserScore(c *gin.Context) {
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
	score, scoreErr := functions.GetUserAvgScore(userID)
	if scoreErr != nil {
		functions.LogError(functions.GetFile(), scoreErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"score": score})
}

func GetTopUsers(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	_, sessionStatus, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(sessionStatus, gin.H{"message": sessionErr.Error()})
		return
	}

	topUsers, topUsersErr := functions.GetTopUsers()
	if topUsersErr != nil {
		functions.LogError(functions.GetFile(), topUsersErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"topUsers": topUsers})
}
