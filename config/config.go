package config

import (
	"fmt"
	"os"
	"strconv"
)

var serviceVersion = "local"

const (
	HTTP_PORT      = "HTTP_PORT"
	LEADING_ZEROES = "LEADING_ZEROES"
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

func GetEnvInt(key string) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, fmt.Errorf("Env Val Not Found")
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("unable to convert env value to int")
	}
	return intVal, nil
}
