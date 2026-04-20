package types

import "time"

type Message struct {
	ID        int64     `json:"id,omitempty"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	TopicID   int64     `json:"topic_id"`
}
