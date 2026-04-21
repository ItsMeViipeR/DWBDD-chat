package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

	authGroup := r.Group("/api")

	authGroup.Use(AuthMiddleware())
	{
		authGroup.DELETE("/messages/:id", DeleteMessageHandler)
		authGroup.POST("/messages", CreateMessageHandler)
		authGroup.POST("/change_email", CreateEmailHandler)
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

	r.POST("/api/login", func(c *gin.Context) {
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
	})

	r.POST("/api/topics", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) < 8 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
			return
		}

		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			return
		}

		userID := int64(claims["user_id"].(float64))

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
