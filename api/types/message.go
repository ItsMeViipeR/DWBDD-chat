package types

import "time"

type Message struct {
	ID        int       `gorm:"primaryKey"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
