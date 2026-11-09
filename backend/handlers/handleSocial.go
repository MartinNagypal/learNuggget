package handlers

import (
	"backend/functions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserScore(c *gin.Context) {
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
	score, scoreEr := functions.GetUserAvgScore(userID)
	if scoreEr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": scoreEr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"score": score})
}

func GetTopUsers(c *gin.Context) {
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

	topUsers, topUsersErr := functions.GetTopUsers()
	if topUsersErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": topUsersErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"topUsers": topUsers})
}
