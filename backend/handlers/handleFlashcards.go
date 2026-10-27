package handlers

import (
	"backend/functions"
	"log"

	"github.com/gin-gonic/gin"
)

func AddNewCard(c *gin.Context) {
	var newCard functions.Card
	err := c.ShouldBindJSON(&newCard)
	if err != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(err)
		c.JSON(500, gin.H{"message": err})
	}

	dbErr := functions.AddNewCard(newCard.Question, newCard.Answer, newCard.CollectionID)
	if dbErr != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(err)
		c.JSON(500, gin.H{"message": dbErr})
	}

	c.JSON(200, gin.H{"message": "Card added succesfully"})
}
