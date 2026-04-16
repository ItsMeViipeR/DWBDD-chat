package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

type Message struct {
	ID        int       `gorm:"primaryKey"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/api/register", func(c *gin.Context) {

	})

	r.POST("/api/login", func(c *gin.Context) {

	})

	r.GET("/api/messages", func(c *gin.Context) {

	})

	r.POST("/api/messages", func(c *gin.Context) {

	})

	r.Run(":8080")
}
