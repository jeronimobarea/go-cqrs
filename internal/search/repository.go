package search

import (
	"context"
)

type Repository interface {
	Index(ctx context.Context, idx, id string, data any) error
	Search(ctx context.Context, idx, query string, fields []string) ([]any, error)
}
