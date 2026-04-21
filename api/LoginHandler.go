package main

import (
	"net/http"
	"os"
	"time"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données manquantes"})
		return
	}

	var users []types.User
	_, err := client.From("users").Select("*", "exact", false).Eq("email", input.Email).ExecuteTo(&users)

	if err != nil || len(users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email inconnu"})
		return
	}

	dbUser := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  dbUser.ID,
		"username": dbUser.Username,
		"email":    dbUser.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la signature"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Connexion réussie",
		"token":   tokenString,
		"user":    dbUser.Username,
	})
}
