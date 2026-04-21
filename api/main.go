package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/postgrest-go"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func main() {
	initDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/login", LoginHandler)

	authGroup := r.Group("/api")

	authGroup.Use(AuthMiddleware())
	{
		authGroup.DELETE("/messages/:id", DeleteMessageHandler)
		authGroup.POST("/messages", CreateMessageHandler)
		authGroup.POST("/change_email", CreateEmailHandler)
		authGroup.POST("/topics", CreateTopicHandler)
	}

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
	})

	r.GET("/api/messages", func(c *gin.Context) {
		var input types.GetMessagesInput
		if err := c.ShouldBindQuery(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var topic []map[string]any
		_, topicErr := client.From("topics").
			Select("id", "exact", false).
			Eq("id", fmt.Sprintf("%d", input.TopicID)).
			ExecuteTo(&topic)

		if topicErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification du sujet"})
			return
		}

		if len(topic) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Le sujet demandé n'existe pas"})
			return
		}

		var result []types.Message
		_, dbErr := client.From("messages").
			Select("*, user:users(id, username)", "exact", false).
			Eq("topic_id", fmt.Sprintf("%d", input.TopicID)).
			Order("created_at", &postgrest.OrderOpts{Ascending: true}).
			ExecuteTo(&result)

		if dbErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": result})
	})

	r.GET("/api/topics", func(c *gin.Context) {
		var result []map[string]any
		_, dbErr := client.From("topics").Select("*", "exact", false).ExecuteTo(&result)

		if dbErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"topics": result})
	})

	r.Run(":8080")
}
