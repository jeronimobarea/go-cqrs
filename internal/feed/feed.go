package feed

import (
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID           uuid.UUID
	Title        string
	Description  string
	CreationDate time.Time
}
