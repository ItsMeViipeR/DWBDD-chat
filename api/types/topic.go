package types

import "time"

type Topic struct {
	ID          int64     `json:"id,omiempty"`
	CreatedAt   time.Time `json:"created_at,omitzero"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
}
