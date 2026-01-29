package postgres

import "github.com/MihirSahani/Project-27/internal"

type PostgresConfig struct {
	address  string
	MaxIdleConns int
	MaxOpenConns int
	MaxIdleTime  string
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		address: internal.GetEnvAsString("POSTGRES_ADDRESS", ""),
		MaxIdleConns: 10,
		MaxOpenConns: 25,
		MaxIdleTime:  "15m",
	}
}