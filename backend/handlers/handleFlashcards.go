package handlers

import (
	"backend/functions"
	"log"

	"github.com/gin-gonic/gin"
)

type newCardType struct {
	question     string
	answer       string
	collectionID int
}

func AddNewCard(c *gin.Context) {
	var newCard newCardType
	err := c.ShouldBindJSON(&newCard)
	if err != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(err)
		c.JSON(500, gin.H{"message": err})
	}

	dbErr := functions.AddNewCard(newCard.question, newCard.answer, newCard.collectionID)
	if dbErr != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(err)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Card added succesfully"})
}
