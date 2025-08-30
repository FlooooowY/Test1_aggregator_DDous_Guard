package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		expectedTCP     string
		expectedHTTP    string
		expectedTimeout time.Duration
	}{
		{
			name:            "default values when no env vars set",
			envVars:         map[string]string{},
			expectedTCP:     "8080",
			expectedHTTP:    "8081",
			expectedTimeout: 30 * time.Second,
		},
		{
			name: "custom values from env vars",
			envVars: map[string]string{
				"TCP_PORT":         "9090",
				"HTTP_PORT":        "9091",
				"SHUTDOWN_TIMEOUT": "60s",
			},
			expectedTCP:     "9090",
			expectedHTTP:    "9091",
			expectedTimeout: 60 * time.Second,
		},
		{
			name: "partial env vars",
			envVars: map[string]string{
				"TCP_PORT": "5000",
			},
			expectedTCP:     "5000",
			expectedHTTP:    "8081",
			expectedTimeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения для теста
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Очищаем переменные после теста
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg := Load()

			if cfg.TCP.Port != tt.expectedTCP {
				t.Errorf("Expected TCP port %s, got %s", tt.expectedTCP, cfg.TCP.Port)
			}

			if cfg.HTTP.Port != tt.expectedHTTP {
				t.Errorf("Expected HTTP port %s, got %s", tt.expectedHTTP, cfg.HTTP.Port)
			}

			if cfg.App.ShutdownTimeout != tt.expectedTimeout {
				t.Errorf("Expected shutdown timeout %v, got %v", tt.expectedTimeout, cfg.App.ShutdownTimeout)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns env value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "returns default when env not set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "returns default when env is empty string",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "   ",
			expected:     "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		expected     time.Duration
	}{
		{
			name:         "returns parsed duration when valid",
			key:          "TIMEOUT",
			defaultValue: 30 * time.Second,
			envValue:     "60s",
			expected:     60 * time.Second,
		},
		{
			name:         "returns default when env not set",
			key:          "TIMEOUT",
			defaultValue: 30 * time.Second,
			envValue:     "",
			expected:     30 * time.Second,
		},
		{
			name:         "returns default when env is invalid",
			key:          "TIMEOUT",
			defaultValue: 30 * time.Second,
			envValue:     "invalid",
			expected:     30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvAsDuration(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "returns parsed int when valid",
			key:          "PORT",
			defaultValue: 8080,
			envValue:     "9090",
			expected:     9090,
		},
		{
			name:         "returns default when env not set",
			key:          "PORT",
			defaultValue: 8080,
			envValue:     "",
			expected:     8080,
		},
		{
			name:         "returns default when env is invalid",
			key:          "PORT",
			defaultValue: 8080,
			envValue:     "not_a_number",
			expected:     8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvAsInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
