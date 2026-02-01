package internal

import (
	"fmt"
	"os"
)

func GetEnvAsString(key, defaultValue string) string {
	value, found:= os.LookupEnv(key)
	if !found {
		fmt.Printf("Environment variable not found\n")
		return defaultValue
	} else {
		return value
	}
}