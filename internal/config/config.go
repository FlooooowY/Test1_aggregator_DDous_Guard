package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	TCP  TCPConfig
	HTTP HTTPConfig
	App  AppConfig
}

type TCPConfig struct {
	Port string
}

type HTTPConfig struct {
	Port string
}

type AppConfig struct {
	ShutdownTimeout time.Duration
}

func Load() *Config {
	return &Config{
		TCP: TCPConfig{
			Port: getEnv("TCP_PORT", "8080"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8081"),
		},
		App: AppConfig{
			ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
