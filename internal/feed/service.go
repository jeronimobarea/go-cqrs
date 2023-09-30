package feed

import (
	"context"
	"fmt"

	"github.com/jeronimobarea/go-cqrs/internal/eventstorer"
)

type Service interface {
	Create(ctx context.Context, cmd CreateFeedCmd) error
}

type service struct {
	repo Repository
	nats eventstorer.EventStorer
}

func NewService(repo Repository, nats eventstorer.EventStorer) *service {
	return &service{
		repo: repo,
		nats: nats,
	}
}

func (s *service) Create(ctx context.Context, cmd CreateFeedCmd) error {
	feed := newFeed(cmd.Title, cmd.Description)
	err := s.repo.Insert(ctx, feed)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	msg := CreatedFeedMessage{
		ID:           feed.ID,
		Title:        feed.Title,
		Description:  feed.Description,
		CreationDate: feed.CreationDate,
	}
	err = s.nats.Publish(ctx, msg)
	if err != nil {
		return fmt.Errorf("error publishing created feed msg: %w", err)
	}
	return nil
}
