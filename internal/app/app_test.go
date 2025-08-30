package app

import (
	"aggregator/internal/config"
	"testing"
	"time"
)

func TestNewApp(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectedNil bool
	}{
		{
			name: "creates app successfully with valid config",
			config: &config.Config{
				TCP:  config.TCPConfig{Port: "8080"},
				HTTP: config.HTTPConfig{Port: "8081"},
				App:  config.AppConfig{ShutdownTimeout: 30 * time.Second},
			},
			expectedNil: false,
		},
		{
			name:        "creates app with nil config",
			config:      nil,
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.config)

			if tt.expectedNil && app != nil {
				t.Error("Expected app to be nil")
			}
			if !tt.expectedNil && app == nil {
				t.Error("Expected app to be created")
			}

			if app != nil {
				if app.config != tt.config {
					t.Error("Expected config to be set correctly")
				}
				if app.server == nil {
					t.Error("Expected server to be initialized")
				}
				if app.quit == nil {
					t.Error("Expected quit channel to be initialized")
				}
			}
		})
	}
}

func TestApp_LogStartupInfo(t *testing.T) {
	tests := []struct {
		name     string
		tcpPort  string
		httpPort string
	}{
		{
			name:     "logs startup info with default ports",
			tcpPort:  "8080",
			httpPort: "8081",
		},
		{
			name:     "logs startup info with custom ports",
			tcpPort:  "9090",
			httpPort: "9091",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				TCP:  config.TCPConfig{Port: tt.tcpPort},
				HTTP: config.HTTPConfig{Port: tt.httpPort},
				App:  config.AppConfig{ShutdownTimeout: 30 * time.Second},
			}

			app := NewApp(cfg)
			// logStartupInfo - это приватный метод, но мы можем проверить, что конфигурация установлена правильно
			if app.config.TCP.Port != tt.tcpPort {
				t.Errorf("Expected TCP port %s, got %s", tt.tcpPort, app.config.TCP.Port)
			}
			if app.config.HTTP.Port != tt.httpPort {
				t.Errorf("Expected HTTP port %s, got %s", tt.httpPort, app.config.HTTP.Port)
			}
		})
	}
}

func TestApp_ServerInitialization(t *testing.T) {
	tests := []struct {
		name         string
		config       *config.Config
		expectedTCP  bool
		expectedHTTP bool
	}{
		{
			name: "server components are properly initialized",
			config: &config.Config{
				TCP:  config.TCPConfig{Port: "8080"},
				HTTP: config.HTTPConfig{Port: "8081"},
				App:  config.AppConfig{ShutdownTimeout: 30 * time.Second},
			},
			expectedTCP:  true,
			expectedHTTP: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.config)

			if tt.expectedTCP && app.server.tcpHandler == nil {
				t.Error("Expected TCP handler to be initialized")
			}
			if tt.expectedHTTP && app.server.httpHandler == nil {
				t.Error("Expected HTTP handler to be initialized")
			}
		})
	}
}
