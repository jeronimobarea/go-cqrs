package feed

import "context"

type Repository interface {
	Insert(ctx context.Context, feed *Feed) error
	List(ctx context.Context) ([]*Feed, error)
}
