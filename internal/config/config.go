package config

import "os"

type Config struct {
	PsqlDSN string
}

func LoadConfig() (*Config, error) {
	return &Config{
		PsqlDSN: os.Getenv("PSQL_DSN"),
	}, nil
}
