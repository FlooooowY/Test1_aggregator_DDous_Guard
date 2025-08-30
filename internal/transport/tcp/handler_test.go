package tcp

import (
	"net"
	"testing"
)

type MockStatsUseCase struct {
	messages []string
}

func (m *MockStatsUseCase) ProcessMessage(content string) {
	m.messages = append(m.messages, content)
}

func (m *MockStatsUseCase) GetStats() (int, int) {
	return len(m.messages), 0
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name        string
		expectedNil bool
	}{
		{
			name:        "creates handler successfully",
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockStatsUseCase{}
			handler := NewHandler(mockUseCase)

			if tt.expectedNil && handler != nil {
				t.Error("Expected handler to be nil")
			}
			if !tt.expectedNil && handler == nil {
				t.Error("Expected handler to be created")
			}

			if handler != nil && handler.statsUseCase == nil {
				t.Error("Expected stats use case to be set")
			}
		})
	}
}

func TestHandler_HandleConnection(t *testing.T) {
	tests := []struct {
		name          string
		messages      []string
		expectedCount int
	}{
		{
			name:          "handles single message",
			messages:      []string{"Hello world"},
			expectedCount: 1,
		},
		{
			name:          "handles multiple messages",
			messages:      []string{"First message", "Second message", "Third message"},
			expectedCount: 3,
		},
		{
			name:          "handles empty message",
			messages:      []string{""},
			expectedCount: 1,
		},
		{
			name:          "handles no messages",
			messages:      []string{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockStatsUseCase{}
			handler := NewHandler(mockUseCase)

			server, client := net.Pipe()
			defer server.Close()
			defer client.Close()

			go func() {
				for _, msg := range tt.messages {
					client.Write([]byte(msg + "\n"))
				}
				client.Close()
			}()

			handler.HandleConnection(server)

			if len(mockUseCase.messages) != tt.expectedCount {
				t.Errorf("Expected %d messages, got %d", tt.expectedCount, len(mockUseCase.messages))
			}

			for i, expected := range tt.messages {
				if i < len(mockUseCase.messages) && mockUseCase.messages[i] != expected {
					t.Errorf("Expected message %d to be '%s', got '%s'", i, expected, mockUseCase.messages[i])
				}
			}
		})
	}
}

func TestHandler_HandleConnection_ConnectionClosed(t *testing.T) {
	tests := []struct {
		name             string
		closeImmediately bool
		expectedCount    int
	}{
		{
			name:             "handles immediate connection close",
			closeImmediately: true,
			expectedCount:    0,
		},
		{
			name:             "handles normal connection",
			closeImmediately: false,
			expectedCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockStatsUseCase{}
			handler := NewHandler(mockUseCase)

			server, client := net.Pipe()
			defer server.Close()
			defer client.Close()

			if tt.closeImmediately {
				client.Close()
			} else {
				go func() {
					client.Write([]byte("test message\n"))
					client.Close()
				}()
			}

			handler.HandleConnection(server)

			if len(mockUseCase.messages) != tt.expectedCount {
				t.Errorf("Expected %d messages, got %d", tt.expectedCount, len(mockUseCase.messages))
			}
		})
	}
}
