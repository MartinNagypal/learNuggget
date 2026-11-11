package handlers

import (
	"backend/functions"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddNewCard(c *gin.Context) {
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

	var newCard functions.Card
	err := c.ShouldBindJSON(&newCard)
	if err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if !requireCollectionOwner(c, newCard.CollectionID, userID) {
		return
	}

	dbErr := functions.AddNewCard(newCard.Question, newCard.Answer, newCard.CollectionID)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Card added succesfully"})
}

func GetCards(c *gin.Context) {
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

	id, convErr := strconv.Atoi(c.Query("collectionID"))
	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(400, gin.H{"message": "No collectionID given"})
		return
	}
	if !requireCollectionOwner(c, id, userID) {
		return
	}

	cards, fetchErr := functions.FetchCards(id)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(500, gin.H{"message": fetchErr.Error()})
		return
	}

	c.JSON(200, cards)
}

func AddBulkCards(c *gin.Context) {
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

	id, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No collection id given"})
		return
	}
	if !requireCollectionOwner(c, id, userID) {
		return
	}
	var body struct {
		Bulk string `json:"bulk"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	failed := false
	bulkCardsCol := strings.Split(body.Bulk, ";;")
	for _, element := range bulkCardsCol {
		if element == "" {
			continue
		}
		sepElement := strings.SplitN(element, "@@", 2)
		if len(sepElement) != 2 {
			continue
		}
		addErr := functions.AddNewCard(sepElement[0], sepElement[1], id)
		if addErr != nil {
			functions.LogError(functions.GetFile(), addErr)
			failed = true
		}
	}

	if failed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "One or more cards could not be added..."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cards added succesfully"})
}

func DeleteCard(c *gin.Context) {
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

	id, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No card id given"})
		return
	}
	err := functions.DeleteCard(id, userID)
	if errors.Is(err, functions.ErrCardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted Succesfully"})
}

func MarkCard(c *gin.Context) {
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

	var req functions.MarkCardRequest
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": bindErr.Error()})
		return
	}

	markErr := functions.MarkCard(req.CardID, req.Marked, userID)
	if markErr != nil {
		functions.LogError(functions.GetFile(), markErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": markErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card updated"})
}

func GetMarkedCards(c *gin.Context) {
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

	cards, fetchErr := functions.FetchMarkedCards(userID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": fetchErr.Error()})
		return
	}
	c.JSON(http.StatusOK, cards)
}
