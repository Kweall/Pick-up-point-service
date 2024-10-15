package config

import "os"

type Config struct {
	PsqlDSN   string
	GrpcHost  string
	HttpHost  string
	AdminHost string
}

func LoadConfig() (*Config, error) {
	return &Config{
		PsqlDSN:   os.Getenv("PSQL_DSN"),
		GrpcHost:  os.Getenv("grpcHost"),
		HttpHost:  os.Getenv("httpHost"),
		AdminHost: os.Getenv("adminHost"),
	}, nil
}
