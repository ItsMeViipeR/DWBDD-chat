package main

import (
	"fmt"
	"net/http"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var input types.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 0)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser := types.User{
		Username: input.Name,
		Email:    input.Email,
		Password: string(hash),
	}

	var result []types.User

	_, err = client.From("users").Insert(newUser, false, "", "", "").ExecuteTo(&result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la sauvegarde"})
		fmt.Println("Erreur lors de la sauvegarde", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Utilisateur créé avec succès",
		"user":    result[0].Username,
	})
}
