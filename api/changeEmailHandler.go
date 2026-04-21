package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateEmailHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'email invalide"})
		fmt.Println("c: ", c)
		return
	}

	var updatedUsers []types.User

	_, err := client.From("users").Update(map[string]any{"email": input.Email}, "", "").Eq("id", fmt.Sprintf("%d", userID)).ExecuteTo(&updatedUsers)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'email"})
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   updatedUsers[0].Email,
	})
	newTokenString, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Email mis à jour avec succès",
		"new_email": updatedUsers[0].Email,
		"token":     newTokenString,
	})
}
