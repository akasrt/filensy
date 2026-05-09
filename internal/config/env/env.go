package env

import "os"

const (
	// a list of env key constants
	DSN        = "DSN"
	FileRoot   = "FILE_ROOT"
	ServerPort = "SERVER_PORT"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}
