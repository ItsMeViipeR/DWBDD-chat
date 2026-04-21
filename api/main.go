package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	r.POST("/api/register", RegisterHandler)
	r.POST("/api/login", LoginHandler)
	r.GET("/api/messages", GetMessagesHandler)

	authGroup := r.Group("/api")

	authGroup.Use(AuthMiddleware())
	{
		authGroup.DELETE("/messages/:id", DeleteMessageHandler)
		authGroup.POST("/messages", CreateMessageHandler)
		authGroup.POST("/change_email", CreateEmailHandler)
		authGroup.POST("/topics", CreateTopicHandler)
	}

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
