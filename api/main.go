package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ItsMeViipeR/DWBDD-chat/api/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/postgrest-go"
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

	authGroup := r.Group("/api")

	authGroup.Use(AuthMiddleware())
	{
		authGroup.DELETE("/messages/:id", DeleteMessageHandler)
		authGroup.POST("/messages", CreateMessageHandler)
		authGroup.POST("/change_email", CreateEmailHandler)
		authGroup.POST("/topics", CreateTopicHandler)
	}

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
