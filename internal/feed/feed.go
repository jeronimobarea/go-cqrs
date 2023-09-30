package feed

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID           uuid.UUID
	Title        string
	Description  string
	CreationDate time.Time
}

func newFeed(title, description string) Feed {
	return Feed{
		ID:           uuid.New(),
		Title:        title,
		Description:  description,
		CreationDate: time.Now(),
	}
}

type CreateFeedCmd struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func ParseSearchHits(hits []any) ([]*Feed, error) {
	var feeds []*Feed
	for _, hit := range hits {
		feed, err := parseSearchHit(hit)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func parseSearchHit(hit any) (*Feed, error) {
	source := hit.(map[string]any)["_source"]
	m, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	var feed Feed
	err = json.Unmarshal(m, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}
