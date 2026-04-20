package types

import "time"

type Topic struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
}
