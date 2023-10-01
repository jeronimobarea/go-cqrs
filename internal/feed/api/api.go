package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"

	"github.com/jeronimobarea/go-cqrs/internal/config"
	eventStorer "github.com/jeronimobarea/go-cqrs/internal/eventstorer/nats"
	"github.com/jeronimobarea/go-cqrs/internal/feed"
	feedHandlers "github.com/jeronimobarea/go-cqrs/internal/feed/handlers"
	feedRepo "github.com/jeronimobarea/go-cqrs/internal/feed/repository"
)

var cfg config.Config

func init() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
}

func Run() {
	var feedSvc feed.Service
	{
		dbUrl := buildDBUrl()
		db, err := sql.Open("postgres", dbUrl)
		if err != nil {
			panic(err)
		}
		repo := feedRepo.NewRepository(db)

		natsUrl := buildNatsUrl()
		conn, err := nats.Connect(natsUrl)
		if err != nil {
			panic(err)
		}
		eventStorer := eventStorer.NewEventStorer(conn)
		feedSvc = feed.NewService(repo, eventStorer)
	}

	var r *chi.Mux
	{
		r = chi.NewRouter()

		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))
	}

	feedHandlers.RegisterRoutes(r, feedSvc)

	http.ListenAndServe(":3000", r)
}

func buildDBUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@postgres/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DB,
	)
}

func buildNatsUrl() string {
	return fmt.Sprintf("nats://%s", cfg.Nats.Address)
}
