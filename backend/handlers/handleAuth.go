package handlers

import (
	"backend/functions"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Status codes used by all handlers:
//
//	200 success
//	400 invalid input (broken JSON, missing/invalid id, validation)
//	401 not logged in, session not found/expired, wrong login data
//	403 logged in but not allowed (no admin, protected root user)
//	404 resource does not exist or does not belong to the user
//	409 username already taken
//	429 too many login attempts
//	500 internal error: LogError(GetFile(), err) + "Internal server error"
//
// Only 500 errors are logged; 4xx are expected user errors.

const isProduction = false

func RegisterUser(c *gin.Context) {
	var newUser functions.UserElement
	bindErr := c.ShouldBindJSON(&newUser)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	username, usernameErr := functions.NormalizeUsername(newUser.Username)
	if usernameErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": usernameErr.Error()})
		return
	}

	pwErr := functions.ValidatePw(newUser.Password)
	if pwErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": pwErr.Error()})
		return
	}

	exists, existsErr := functions.CheckIfUserExists(username)
	if existsErr != nil {
		functions.LogError(functions.GetFile(), existsErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "User is already existing"})
		return
	}
	sessionToken, createErr := functions.CreateUser(username, newUser.Password)
	if createErr != nil {
		functions.LogError(functions.GetFile(), createErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	// SameSite=Strict: the browser never sends the cookie with requests started by another site
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("sessionToken", sessionToken, functions.SessionLifetimeSeconds, "/", "", isProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}

func LoginUser(c *gin.Context) {
	var user functions.UserElement
	bindErr := c.ShouldBindJSON(&user)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Login returns 401 (wrong data), 429 (locked) or 500 (internal) on errors
	sessionToken, httpStatus, loginErr := functions.Login(strings.TrimSpace(user.Username), user.Password)
	if loginErr != nil {
		if httpStatus == http.StatusInternalServerError {
			functions.LogError(functions.GetFile(), loginErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
		c.JSON(httpStatus, gin.H{"message": loginErr.Error()})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("sessionToken", sessionToken, functions.SessionLifetimeSeconds, "/", "", isProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}

func VerifySessionValid(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Valid session"})
}

func Logout(c *gin.Context) {
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
	logoutErr := functions.FuncLogout(sessionToken)
	if logoutErr != nil {
		functions.LogError(functions.GetFile(), logoutErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	// Max-Age -1: the browser deletes the cookie
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("sessionToken", "", -1, "/", "", isProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func FetchUsername(c *gin.Context) {
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
	username, usernameErr := functions.GetUsername(userID)
	if usernameErr != nil {
		functions.LogError(functions.GetFile(), usernameErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	root, rootErr := functions.IsRoot(userID)
	if rootErr != nil {
		functions.LogError(functions.GetFile(), rootErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username, "root": root})
}

func DelAccount(c *gin.Context) {
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
	delErr := functions.DeleteAccount(userID)
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
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("sessionToken", "", -1, "/", "", isProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted account"})
}

func LogoutAllSession(c *gin.Context) {
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
	logoutErr := functions.LogoutAllSessions(userID)
	if logoutErr != nil {
		functions.LogError(functions.GetFile(), logoutErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("sessionToken", "", -1, "/", "", isProduction, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out from all sessions"})
}
