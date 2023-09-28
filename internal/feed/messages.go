package feed

import (
	"time"

	"github.com/google/uuid"
)

type CreatedFeedMessage struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
}

func (m CreatedFeedMessage) Type() string {
	return "created_feed"
}
