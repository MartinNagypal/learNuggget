package handlers

import (
	"backend/functions"
	"net/http"

	"github.com/gin-gonic/gin"
)

const isProduction = false

func RegisterUser(c *gin.Context) {
	var newUser functions.RegisterUserElement
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
