package pubsub

import (
	"context"
)

type EventStorer interface {
	Close()
	Publish(ctx context.Context, msg Message) error
	Subscribe(ctx context.Context) (<-chan Message, error)
	OnCreate(f func(Message)) error
}
