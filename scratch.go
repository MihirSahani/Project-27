package main

import (
	"fmt"

	"github.com/MihirSahani/Project-27/storage/cache/redis"
)

func main() {
	cacheManager := redis.NewRedisCacheManager()
	err := cacheManager.Ping()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		panic(err)
	}
	fmt.Printf("Ping success")
}