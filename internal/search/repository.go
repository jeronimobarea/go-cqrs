package search

import (
	"context"
	"io"
)

type Repository interface {
	Index(ctx context.Context, idx, id string, data any) error
	Search(ctx context.Context, idx, query string, fields []string) (io.ReadCloser, error)
}
