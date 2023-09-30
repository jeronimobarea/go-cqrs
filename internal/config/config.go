package config

type (
	postgres struct {
		DB       string `envconfig:"POSTGRES_DB"`
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
	}

	nats struct {
		Address string `envconfig:"NATS_ADDRESS"`
	}

	Config struct {
		Postgres postgres
		Nats     nats
	}
)
