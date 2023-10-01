package eventstorer

import (
	"context"
)

type EventStorer interface {
	Publish(ctx context.Context, msg Message) error
	Subscribe(ctx context.Context) (<-chan Message, error)
	OnCreate(f func(Message)) error
}
