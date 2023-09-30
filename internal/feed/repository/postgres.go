package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/jeronimobarea/go-cqrs/internal/feed"
)

var _ feed.Repository = &repository{}

const (
	tableName   = "feeds"
	tableFields = "id, uuid, title, description, creation_date"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Insert(ctx context.Context, feed feed.Feed) error {
	q := `INSERT INTO ` + tableName + `(uuid, title, description, creation_date)` +
		`VALUES ($1, $2, $3, $4)`

	args := []interface{}{
		feed.ID,
		feed.Title,
		feed.Description,
		feed.CreationDate,
	}

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *repository) List(ctx context.Context) ([]*feed.Feed, error) {
	q := `SELECT uuid, title, description, creation_Date FROM ` + tableName

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*feed.Feed
	for rows.Next() {
		var feed *feed.Feed
		if err := rows.Scan(
			&feed.ID,
			&feed.Title,
			&feed.Description,
			&feed.CreationDate,
		); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}
