package usecase

import (
	"aggregator/internal/domain"
	"testing"
)

// MockStatsRepository для тестирования
type MockStatsRepository struct {
	messages    []*domain.Message
	totalLines  int
	uniqueWords int
}

func (m *MockStatsRepository) AddMessage(message *domain.Message) {
	m.messages = append(m.messages, message)
	m.totalLines++
	// Упрощенная логика для мока
	m.uniqueWords = len(message.Words)
}

func (m *MockStatsRepository) GetStats() (int, int) {
	return m.totalLines, m.uniqueWords
}

func TestNewStatsUseCase(t *testing.T) {
	tests := []struct {
		name        string
		expectedNil bool
	}{
		{
			name:        "new use case should be created successfully",
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStatsRepository{}
			uc := NewStatsUseCase(mockRepo)

			if tt.expectedNil && uc != nil {
				t.Error("Expected use case to be nil")
			}
			if !tt.expectedNil && uc == nil {
				t.Error("Expected use case to be created")
			}
		})
	}
}

func TestStatsUseCase_ProcessMessage(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		expectedCount   int
		expectedContent string
	}{
		{
			name:            "process simple message",
			message:         "Hello world",
			expectedCount:   1,
			expectedContent: "Hello world",
		},
		{
			name:            "process message with punctuation",
			message:         "Hello, world! How are you?",
			expectedCount:   1,
			expectedContent: "Hello, world! How are you?",
		},
		{
			name:            "process empty message",
			message:         "",
			expectedCount:   0,
			expectedContent: "",
		},
		{
			name:            "process message with spaces only",
			message:         "   ",
			expectedCount:   0,
			expectedContent: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStatsRepository{}
			uc := NewStatsUseCase(mockRepo)

			uc.ProcessMessage(tt.message)

			if len(mockRepo.messages) != tt.expectedCount {
				t.Errorf("Expected %d messages, got %d", tt.expectedCount, len(mockRepo.messages))
			}

			if tt.expectedCount > 0 && mockRepo.messages[0].Content != tt.expectedContent {
				t.Errorf("Expected '%s', got '%s'", tt.expectedContent, mockRepo.messages[0].Content)
			}
		})
	}
}

func TestStatsUseCase_GetStats(t *testing.T) {
	tests := []struct {
		name          string
		totalLines    int
		uniqueWords   int
		expectedLines int
		expectedWords int
	}{
		{
			name:          "get stats with zero values",
			totalLines:    0,
			uniqueWords:   0,
			expectedLines: 0,
			expectedWords: 0,
		},
		{
			name:          "get stats with positive values",
			totalLines:    5,
			uniqueWords:   10,
			expectedLines: 5,
			expectedWords: 10,
		},
		{
			name:          "get stats with large values",
			totalLines:    1000,
			uniqueWords:   500,
			expectedLines: 1000,
			expectedWords: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStatsRepository{
				totalLines:  tt.totalLines,
				uniqueWords: tt.uniqueWords,
			}
			uc := NewStatsUseCase(mockRepo)

			totalLines, uniqueWords := uc.GetStats()

			if totalLines != tt.expectedLines {
				t.Errorf("Expected %d total lines, got %d", tt.expectedLines, totalLines)
			}
			if uniqueWords != tt.expectedWords {
				t.Errorf("Expected %d unique words, got %d", tt.expectedWords, uniqueWords)
			}
		})
	}
}
