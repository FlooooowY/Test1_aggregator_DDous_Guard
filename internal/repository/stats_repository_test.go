package repository

import (
	"aggregator/internal/domain"
	"sync"
	"testing"
)

func TestNewInMemoryStatsRepository(t *testing.T) {
	tests := []struct {
		name          string
		expectedLines int
		expectedWords int
	}{
		{
			name:          "new repository should be initialized with zero values",
			expectedLines: 0,
			expectedWords: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryStatsRepository()

			totalLines, uniqueWords := repo.GetStats()
			if totalLines != tt.expectedLines {
				t.Errorf("Expected %d total lines, got %d", tt.expectedLines, totalLines)
			}
			if uniqueWords != tt.expectedWords {
				t.Errorf("Expected %d unique words, got %d", tt.expectedWords, uniqueWords)
			}
		})
	}
}

func TestInMemoryStatsRepository_AddMessage(t *testing.T) {
	tests := []struct {
		name          string
		messages      []string
		expectedLines []int
		expectedWords []int
	}{
		{
			name:          "add single message",
			messages:      []string{"Hello world"},
			expectedLines: []int{1},
			expectedWords: []int{2},
		},
		{
			name:          "add multiple messages",
			messages:      []string{"Hello world", "Hello again"},
			expectedLines: []int{1, 2},
			expectedWords: []int{2, 3},
		},
		{
			name:          "add message with punctuation",
			messages:      []string{"Hello, world! How are you?"},
			expectedLines: []int{1},
			expectedWords: []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryStatsRepository()

			for i, message := range tt.messages {
				msg := domain.NewMessage(message)
				repo.AddMessage(msg)

				totalLines, uniqueWords := repo.GetStats()
				if totalLines != tt.expectedLines[i] {
					t.Errorf("After message %d: Expected %d lines, got %d", i+1, tt.expectedLines[i], totalLines)
				}
				if uniqueWords != tt.expectedWords[i] {
					t.Errorf("After message %d: Expected %d unique words, got %d", i+1, tt.expectedWords[i], uniqueWords)
				}
			}
		})
	}
}

func TestInMemoryStatsRepository_AddMessage_Empty(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		expectedLines int
		expectedWords int
	}{
		{
			name:          "empty message should count as line but no words",
			message:       "",
			expectedLines: 1,
			expectedWords: 0,
		},
		{
			name:          "message with only spaces should count as line but no words",
			message:       "   ",
			expectedLines: 1,
			expectedWords: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryStatsRepository()

			msg := domain.NewMessage(tt.message)
			repo.AddMessage(msg)

			totalLines, uniqueWords := repo.GetStats()
			if totalLines != tt.expectedLines {
				t.Errorf("Expected %d lines, got %d", tt.expectedLines, totalLines)
			}
			if uniqueWords != tt.expectedWords {
				t.Errorf("Expected %d unique words, got %d", tt.expectedWords, uniqueWords)
			}
		})
	}
}

func TestInMemoryStatsRepository_Concurrent(t *testing.T) {
	tests := []struct {
		name                 string
		numGoroutines        int
		messagesPerGoroutine int
		messageContent       string
		expectedLines        int
		expectedWords        int
	}{
		{
			name:                 "concurrent access with 10 goroutines",
			numGoroutines:        10,
			messagesPerGoroutine: 100,
			messageContent:       "goroutine message",
			expectedLines:        1000,
			expectedWords:        2,
		},
		{
			name:                 "concurrent access with 5 goroutines",
			numGoroutines:        5,
			messagesPerGoroutine: 50,
			messageContent:       "test message",
			expectedLines:        250,
			expectedWords:        2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryStatsRepository()

			var wg sync.WaitGroup

			for i := 0; i < tt.numGoroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for j := 0; j < tt.messagesPerGoroutine; j++ {
						msg := domain.NewMessage(tt.messageContent)
						repo.AddMessage(msg)
					}
				}(i)
			}

			wg.Wait()

			totalLines, uniqueWords := repo.GetStats()
			if totalLines != tt.expectedLines {
				t.Errorf("Expected %d lines, got %d", tt.expectedLines, totalLines)
			}
			if uniqueWords != tt.expectedWords {
				t.Errorf("Expected %d unique words, got %d", tt.expectedWords, uniqueWords)
			}
		})
	}
}
