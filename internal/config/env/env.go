package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	// a list of env key constants
	DSN         = "DSN"
	ServerPort  = "SERVER_PORT"
	FileRoot    = "FILE_ROOT"
	AuthKey     = "AUTH_KEY"
	LogLevel    = "LOG_LEVEL"
	FileCodeKey = "FILE_CODE_KEY"
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
