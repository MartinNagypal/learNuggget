package handlers

import (
	"backend/functions"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddNewCard(c *gin.Context) {
	var newCard functions.Card
	err := c.ShouldBindJSON(&newCard)
	if err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
	}

	dbErr := functions.AddNewCard(newCard.Question, newCard.Answer, newCard.CollectionID)
	if dbErr != nil {
		functions.LogError(functions.GetFile(), dbErr)
		c.JSON(500, gin.H{"message": dbErr.Error()})
	}

	c.JSON(200, gin.H{"message": "Card added succesfully"})
}

func GetCards(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Query("collectionID"))

	cards, fetchErr := functions.FetchCards(id)

	if convErr != nil {
		functions.LogError(functions.GetFile(), convErr)
		c.JSON(400, gin.H{"message": "No collectionID given"})
		return
	}

	if fetchErr != nil {
		functions.LogError(functions.GetFile(), fetchErr)
		c.JSON(500, gin.H{"message": fetchErr.Error()})
		return
	}

	c.JSON(200, cards)
}

func AddBulkCards(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	var body struct {
		Bulk string `json:"bulk"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	type bulkCard struct {
		Question string
		Answer   string
	}

	bulkCardsCol := strings.Split(body.Bulk, ";")
	for _, element := range bulkCardsCol {
		if element == "" {
			continue
		}
		sepElement := strings.Split(element, "#")
		if len(sepElement) != 2 {
			continue
		}
		addErr := functions.AddNewCard(sepElement[0], sepElement[1], id)
		if addErr != nil {
			functions.LogError(functions.GetFile(), addErr)
			c.JSON(500, gin.H{"message": "One or more cards could not be added..."})
			continue
		}
		println(sepElement)
	}
}

func DeleteCard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	err := functions.DeleteCard(id)
	if err != nil {
		functions.LogError(functions.GetFile(), err)
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted Succesfully"})
}
