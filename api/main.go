package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// On ajoute le tag json pour correspondre aux colonnes Supabase
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	ID        int       `gorm:"primaryKey"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

var client *supabase.Client

func initDB() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Erreur chargement .env:", err)
		return
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	client, err = supabase.NewClient(supabaseURL, supabaseKey, nil)

	if err != nil {
		fmt.Println("Erreur init Supabase:", err)
		return
	}
}

func main() {
	initDB()

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/api/register", func(c *gin.Context) {
		var input RegisterInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 0)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newUser := User{
			Username: input.Name,
			Email:    input.Email,
			Password: string(hash),
		}

		var result []User

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
	})

	r.POST("/api/login", func(c *gin.Context) {

	})

	r.GET("/api/messages", func(c *gin.Context) {

	})

	r.POST("/api/messages", func(c *gin.Context) {

	})

	r.Run(":8080")
}
