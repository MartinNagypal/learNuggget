package handlers

import (
	"backend/functions"
	"net/http"

	"github.com/gin-gonic/gin"
)

const isProduction = false

func RegisterUser(c *gin.Context) {
	var newUser functions.UserElement
	bindErr := c.ShouldBindJSON(&newUser)
	if bindErr != nil {
		functions.LogError(functions.GetFile(), bindErr)
		c.JSON(500, gin.H{"message": bindErr.Error()})
	}
	exists, existsErr := functions.CheckIfUserExists(newUser.Username)
	if existsErr != nil {
		functions.LogError(functions.GetFile(), existsErr)
		c.JSON(http.StatusConflict, existsErr.Error())
	}
	if exists == true {
		c.JSON(http.StatusConflict, gin.H{"message": "User is already existing"})
		return
	}
	sessionToken, createErr := functions.CreateUser(newUser.Username, newUser.Password)
	if createErr != nil {
		functions.LogError(functions.GetFile(), createErr)
		c.JSON(http.StatusInternalServerError, createErr.Error())
	}
	c.SetCookie("sessionToken", sessionToken, 432000, "/", "", isProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "user succesfully registered"})
}

func LoginUser(c *gin.Context) {
	var user functions.UserElement
	bindErr := c.ShouldBindJSON(&user)
	if bindErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": bindErr.Error()})
		return
	}

	sessionToken, httpStatus, loginErr := functions.Login(user.Username, user.Password)
	if loginErr != nil {
		c.JSON(httpStatus, gin.H{"message": loginErr.Error()})
		return
	}
	c.SetCookie("sessionToken", sessionToken, 432000, "/", "", isProduction, true)
	c.JSON(httpStatus, gin.H{"message": "Succesfully logged in"})
}

func VerifySessionValid(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	_, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Valid session"})
}

func Logout(c *gin.Context) {
	sessionToken, cookieErr := c.Cookie("sessionToken")
	if cookieErr != nil {
		functions.LogError(functions.GetFile(), cookieErr)
		c.JSON(http.StatusForbidden, gin.H{"message": cookieErr.Error()})
		return
	}
	userID, sessionErr := functions.VerifySession(sessionToken)
	if sessionErr != nil {
		functions.LogError(functions.GetFile(), sessionErr)
		c.JSON(http.StatusForbidden, gin.H{"message": sessionErr.Error()})
		return
	}
	logoutErr := functions.FuncLogout(userID)
	if logoutErr != nil {
		functions.LogError(functions.GetFile(), logoutErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": logoutErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Succesfully logged out"})
}
