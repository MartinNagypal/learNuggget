package handlers

import (
	"backend/functions"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddNewCard(c *gin.Context) {
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

	var newCard functions.Card
	bindErr := c.ShouldBindJSON(&newCard)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	validErr := functions.ValidateCard(newCard.Question, newCard.Answer)
	if validErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": validErr.Error()})
		return
	}
	if !requireCollectionOwner(c, newCard.CollectionID, userID) {
		return
	}

	addErr := functions.AddNewCard(newCard.Question, newCard.Answer, newCard.CollectionID)
	if errors.Is(addErr, functions.ErrTooManyCards) {
		c.JSON(http.StatusBadRequest, gin.H{"message": addErr.Error()})
		return
	}
	if addErr != nil {
		functions.LogError(functions.GetFile(), addErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card added succesfully"})
}

func GetCards(c *gin.Context) {
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

	collectionID, convErr := strconv.Atoi(c.Query("collectionID"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No collection id given"})
		return
	}
	if !requireCollectionOwner(c, collectionID, userID) {
		return
	}

	cards, fetchErr := functions.FetchCards(collectionID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func AddBulkCards(c *gin.Context) {
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
	var body struct {
		Bulk string `json:"bulk"`
	}
	bindErr := c.ShouldBindJSON(&body)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	if !requireCollectionOwner(c, collectionID, userID) {
		return
	}

	// Format: "question@@answer;;question@@answer;;..." – \;; and \@@ for literal separators
	cards := functions.ParseBulkCards(body.Bulk)
	if len(cards) > functions.MaxCardsPerImport {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("At most %d cards per import", functions.MaxCardsPerImport)})
		return
	}
	for _, card := range cards {
		validErr := functions.ValidateCard(card.Question, card.Answer)
		if validErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": validErr.Error()})
			return
		}
	}
	if len(cards) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No valid cards found (format: question@@answer;;)"})
		return
	}

	// One transaction: either all cards are imported or none
	addErr := functions.AddCards(collectionID, cards)
	if errors.Is(addErr, functions.ErrTooManyCards) {
		c.JSON(http.StatusBadRequest, gin.H{"message": addErr.Error()})
		return
	}
	if addErr != nil {
		functions.LogError(functions.GetFile(), addErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d cards added succesfully", len(cards))})
}

func DeleteCard(c *gin.Context) {
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

	cardID, convErr := strconv.Atoi(c.Query("id"))
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No card id given"})
		return
	}
	delErr := functions.DeleteCard(cardID, userID)
	if errors.Is(delErr, functions.ErrCardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": delErr.Error()})
		return
	}
	if delErr != nil {
		functions.LogError(functions.GetFile(), delErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted Succesfully"})
}

func MarkCard(c *gin.Context) {
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

	var req functions.MarkCardRequest
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	markErr := functions.MarkCard(req.CardID, req.Marked, userID)
	if errors.Is(markErr, functions.ErrCardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": markErr.Error()})
		return
	}
	if markErr != nil {
		functions.LogError(functions.GetFile(), markErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card updated"})
}

func GetMarkedCards(c *gin.Context) {
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

	cards, fetchErr := functions.FetchMarkedCards(userID)
	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, cards)
}
