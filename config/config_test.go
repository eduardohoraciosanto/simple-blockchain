package config_test

import (
	"os"
	"testing"

	"github.com/eduardohoraciosanto/simple-blockchain/config"
)

func TestGetVersion(t *testing.T) {
	if config.GetVersion() != "local" {
		t.Fatalf("Unexpected version")
	}
}

func TestGetEnvStringOK(t *testing.T) {
	os.Setenv("TEST_ENV_STR", "TestingString")
	defer os.Unsetenv("TEST_ENV_STR")

	if config.GetEnvString("TEST_ENV_STR", "") != "TestingString" {
		t.Fatalf("Unexpected env value")
	}

}

func TestGetEnvStringDefault(t *testing.T) {
	if config.GetEnvString("TEST_ENV_STR", "default") != "default" {
		t.Fatalf("Unexpected env value")
	}

}

func TestGetPort(t *testing.T) {
	os.Setenv(config.HTTP_PORT, "8001")
	defer os.Unsetenv(config.HTTP_PORT)

	if config.GetPort() != "8001" {
		t.Fatalf("Unexpected Port")
	}
}

func TestGetPortDefault(t *testing.T) {
	os.Unsetenv(config.HTTP_PORT)
	if config.GetPort() != "8080" {
		t.Fatalf("Unexpected Port")
	}
}
