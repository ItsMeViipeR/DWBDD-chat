package main

import (
	"net/http"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-gonic/gin"
)

func CreateTopicHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var input types.CreateTopicInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTopic := types.Topic{
		Name:        input.Name,
		Description: input.Description,
		CreatorID:   userID,
	}

	var result []types.Topic

	_, dbErr := client.From("topics").Insert(newTopic, false, "", "", "").ExecuteTo(&result)

	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Topic créé avec succès", "topic": result[0]})
}
