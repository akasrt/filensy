package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	// a list of env key constants
	DSN           = "DSN"
	ServerAddress = "SERVER_ADDRESS"
	FileRoot      = "FILE_ROOT"
	AuthKey       = "AUTH_KEY"
	LogLevel      = "LOG_LEVEL"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Print(".env file not found")
	}
}
