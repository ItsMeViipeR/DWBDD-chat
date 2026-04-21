package main

import (
	"fmt"
	"net/http"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-gonic/gin"
)

func DeleteMessageHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	messageID := c.Param("id")

	var deletedMessages []types.Message

	_, dbErr := client.From("messages").
		Delete("", "representation").
		Eq("id", messageID).
		Eq("user_id", fmt.Sprintf("%d", userID)).
		ExecuteTo(&deletedMessages)

	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur base de données"})
		return
	}

	if len(deletedMessages) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message non trouvé ou non autorisé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message supprimé"})
}
