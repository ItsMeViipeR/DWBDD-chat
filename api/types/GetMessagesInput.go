package types

type GetMessagesInput struct {
	TopicID int64 `json:"topic_id" form:"topic_id"`
}
