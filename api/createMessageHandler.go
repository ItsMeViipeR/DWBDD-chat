package main

import (
	"net/http"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-gonic/gin"
)

func CreateMessageHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var input types.CreateMessageInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newMessage := types.Message{
		Content: input.Content,
		UserID:  userID,
		TopicID: input.TopicID,
	}

	var result []types.Message

	_, dbErr := client.From("messages").Insert(newMessage, false, "", "", "").ExecuteTo(&result)

	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du message créé"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message créé avec succès", "created_message": result[0]})
}
