package domain

import (
	"sync"
	"testing"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		name           string
		expectedLines  int
		expectedWords  int
		expectedMapNil bool
	}{
		{
			name:           "new stats should be initialized with zero values",
			expectedLines:  0,
			expectedWords:  0,
			expectedMapNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewStats()

			if stats.TotalLines != tt.expectedLines {
				t.Errorf("Expected TotalLines to be %d, got %d", tt.expectedLines, stats.TotalLines)
			}

			if stats.UniqueWords != tt.expectedWords {
				t.Errorf("Expected UniqueWords to be %d, got %d", tt.expectedWords, stats.UniqueWords)
			}

			if tt.expectedMapNil && stats.Words != nil {
				t.Error("Expected words map to be nil")
			}
			if !tt.expectedMapNil && stats.Words == nil {
				t.Error("Expected words map to be initialized")
			}
		})
	}
}

func TestStats_GetStats(t *testing.T) {
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
			totalLines:    2,
			uniqueWords:   3,
			expectedLines: 2,
			expectedWords: 3,
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
			stats := NewStats()

			stats.Mu.Lock()
			stats.TotalLines = tt.totalLines
			stats.UniqueWords = tt.uniqueWords
			stats.Mu.Unlock()

			totalLines, uniqueWords := stats.GetStats()
			if totalLines != tt.expectedLines {
				t.Errorf("Expected %d lines, got %d", tt.expectedLines, totalLines)
			}
			if uniqueWords != tt.expectedWords {
				t.Errorf("Expected %d unique words, got %d", tt.expectedWords, uniqueWords)
			}
		})
	}
}

func TestStats_Concurrent(t *testing.T) {
	tests := []struct {
		name              string
		numGoroutines     int
		callsPerGoroutine int
	}{
		{
			name:              "concurrent access with 10 goroutines",
			numGoroutines:     10,
			callsPerGoroutine: 100,
		},
		{
			name:              "concurrent access with 5 goroutines",
			numGoroutines:     5,
			callsPerGoroutine: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewStats()

			var wg sync.WaitGroup

			for i := 0; i < tt.numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < tt.callsPerGoroutine; j++ {
						stats.GetStats()
					}
				}()
			}

			wg.Wait()
		})
	}
}
