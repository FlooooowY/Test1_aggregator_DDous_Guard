package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockStatsUseCase для тестирования
type MockStatsUseCase struct {
	totalLines  int
	uniqueWords int
}

func (m *MockStatsUseCase) ProcessMessage(content string) {
	// Мок ничего не делает
}

func (m *MockStatsUseCase) GetStats() (int, int) {
	return m.totalLines, m.uniqueWords
}

func TestHandler_GetStats(t *testing.T) {
	tests := []struct {
		name                string
		totalLines          int
		uniqueWords         int
		expectedStatus      int
		expectedContentType string
		expectedTotalLines  int
		expectedUniqueWords int
	}{
		{
			name:                "get stats with positive values",
			totalLines:          2,
			uniqueWords:         4,
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedTotalLines:  2,
			expectedUniqueWords: 4,
		},
		{
			name:                "get stats with zero values",
			totalLines:          0,
			uniqueWords:         0,
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedTotalLines:  0,
			expectedUniqueWords: 0,
		},
		{
			name:                "get stats with large values",
			totalLines:          1000,
			uniqueWords:         500,
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedTotalLines:  1000,
			expectedUniqueWords: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockStatsUseCase{
				totalLines:  tt.totalLines,
				uniqueWords: tt.uniqueWords,
			}
			handler := NewHandler(mockUseCase)

			req, err := http.NewRequest("GET", "/stats", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler.GetStats)

			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if contentType := rr.Header().Get("Content-Type"); contentType != tt.expectedContentType {
				t.Errorf("handler returned wrong content type: got %v want %v", contentType, tt.expectedContentType)
			}

			var response StatsResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to parse JSON response: %v", err)
			}

			if response.TotalLines != tt.expectedTotalLines {
				t.Errorf("Expected total_lines to be %d, got %d", tt.expectedTotalLines, response.TotalLines)
			}
			if response.UniqueWords != tt.expectedUniqueWords {
				t.Errorf("Expected unique_words to be %d, got %d", tt.expectedUniqueWords, response.UniqueWords)
			}
		})
	}
}

func TestHandler_GetStats_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "POST method should return method not allowed",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT method should return method not allowed",
			method:         "PUT",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "DELETE method should return method not allowed",
			method:         "DELETE",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockStatsUseCase{}
			handler := NewHandler(mockUseCase)

			req, err := http.NewRequest(tt.method, "/stats", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler.GetStats)

			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
