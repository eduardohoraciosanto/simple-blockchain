package config

import "os"

var serviceVersion = "local"

const (
	HTTP_PORT = "HTTP_PORT"
)

func GetVersion() string {
	return serviceVersion
}

func GetEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

func GetPort() string {
	return GetEnvString(HTTP_PORT, "8080")
}
