package redis

import (
	"github.com/MihirSahani/Project-27/internal"
)

const (
	DEFAULT_REDIS_DB = 0
)

type RedisConfig struct {
	Address string
	Password string
	Db int
	enabled bool
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Address: internal.GetEnvAsString("REDIS_ADDRESS", ":6379"),
		Password: internal.GetEnvAsString("REDIS_PASSWORD", ""),
		Db: DEFAULT_REDIS_DB,
		enabled: false,
	}
}