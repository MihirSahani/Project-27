package main

import (
	"time"

	"github.com/MihirSahani/Project-27/internal"
)

type AppConfig struct {
	ServerAddress string
	WriteTimeout  time.Duration
	ReadTimeout   time.Duration
	IdleTimeout   time.Duration
	DefaultContextTimeout time.Duration
}

func LoadServerConfig() AppConfig {
	// Load configuration from environment variables or default values
	return AppConfig{
		ServerAddress: internal.GetEnvAsString("BACKEND_SERVER_ADDRESS", ":8080"),
		WriteTimeout:  time.Second * 10,
		ReadTimeout:   time.Second * 10,
		IdleTimeout:   time.Second * 60,
		DefaultContextTimeout: time.Second * 5,
	}
}
