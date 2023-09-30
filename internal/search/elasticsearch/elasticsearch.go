package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/elastic/go-elasticsearch/v8"

	"github.com/jeronimobarea/go-cqrs/internal/search"
)

var _ search.Repository = &repository{}

type repository struct {
	client *elastic.Client
}

func NewRepository(client *elastic.Client) *repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Index(ctx context.Context, idx, id string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = r.client.Index(
		idx,
		bytes.NewReader(b),
		r.client.Index.WithDocumentID(id),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Search(ctx context.Context, idx, query string, fields []string) ([]any, error) {
	searchQuery := NewSearchWithDefaults(query, fields)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(idx),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var rawRes map[string]any
	if err := json.NewDecoder(res.Body).Decode(&rawRes); err != nil {
		return nil, err
	}

	hits := rawRes["hits"].(map[string]any)["hits"].([]any)
	return hits, nil
}
