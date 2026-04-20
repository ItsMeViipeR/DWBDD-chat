package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		var input struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Données manquantes"})
			return
		}

		var users []User
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

	r.GET("/api/user", func(c *gin.Context) {

	})

	r.POST("/api/change_email", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format de token invalide (Bearer manquant)"})
			return
		}

		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Erreur JWT:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des claims"})
			return
		}

		userID := int(claims["user_id"].(float64))

		var input struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'email invalide"})
			fmt.Println("c: ", c)
			return
		}

		var updatedUsers []User

		_, err = client.From("users").Update(map[string]any{"email": input.Email}, "", "").Eq("id", fmt.Sprintf("%d", userID)).ExecuteTo(&updatedUsers)

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
	})

	r.GET("/api/messages", func(c *gin.Context) {

	})

	r.POST("/api/messages", func(c *gin.Context) {

	})

	r.Run(":8080")
}
