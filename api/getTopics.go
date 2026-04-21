package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTopicsHandler(c *gin.Context) {
	var result []map[string]any
	_, dbErr := client.From("topics").Select("*", "exact", false).ExecuteTo(&result)

	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"topics": result})
}
