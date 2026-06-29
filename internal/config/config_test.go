package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("Should use API_PORT when it is a valid number", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "PROD")
		t.Setenv("API_PORT", "9090")

		Load()

		if Port != 9090 {
			t.Errorf("Load should set Port to 9090, got %d", Port)
		}
	})

	t.Run("Should fall back to 8000 when API_PORT is invalid", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "PROD")
		t.Setenv("API_PORT", "not-a-number")

		Load()

		if Port != 8000 {
			t.Errorf("Load should fall back to Port 8000, got %d", Port)
		}
	})

	t.Run("Should build the DB connection string and load the secret key", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "PROD")
		t.Setenv("DB_USER", "user")
		t.Setenv("DB_NAME", "socialnets")
		t.Setenv("DB_PASSWORD", "secret")
		t.Setenv("DB_HOST", "localhost")
		t.Setenv("SECRET_KEY", "my-secret")

		Load()

		expected := "user=user dbname=socialnets password=secret host=localhost sslmode=disable"
		if DbStringConnection != expected {
			t.Errorf("Load should build the connection string.\n expected: %q\n got:      %q", expected, DbStringConnection)
		}
		if string(SecretKey) != "my-secret" {
			t.Errorf("Load should set SecretKey to %q, got %q", "my-secret", string(SecretKey))
		}
	})
}
