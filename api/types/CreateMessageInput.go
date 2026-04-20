package types

type CreateMessageInput struct {
	Content string `json:"content" binding:"required"`
	TopicID int64  `json:"topic_id" binding:"required"`
}
