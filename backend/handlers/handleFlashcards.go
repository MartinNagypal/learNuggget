package handlers

import (
	"backend/functions"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddNewCard(c *gin.Context) {
	var newCard functions.Card
	err := c.ShouldBindJSON(&newCard)
	if err != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(err)
		c.JSON(500, gin.H{"message": err.Error()})
	}

	dbErr := functions.AddNewCard(newCard.Question, newCard.Answer, newCard.CollectionID)
	if dbErr != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(err)
		c.JSON(500, gin.H{"message": dbErr.Error()})
	}

	c.JSON(200, gin.H{"message": "Card added succesfully"})
}

func GetCards(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Query("collectionID"))

	cards, fetchErr := functions.FetchCards(id)

	if convErr != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println("No CollectionID given")
		c.JSON(400, gin.H{"message": "No collectionID given"})
	}

	if fetchErr != nil {
		log.Println("Error in handleFlashcards.go:")
		log.Println(fetchErr)
		c.JSON(500, gin.H{"message": fetchErr.Error()})
	}

	c.JSON(200, cards)

}
