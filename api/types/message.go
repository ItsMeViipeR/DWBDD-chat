package types

import "time"

type Message struct {
	ID        int64
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
