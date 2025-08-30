package domain

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "empty message",
			content:  "",
			expected: []string{},
		},
		{
			name:     "simple message",
			content:  "Hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "message with punctuation",
			content:  "Hello, world! How are you?",
			expected: []string{"hello", "world", "how", "are", "you"},
		},
		{
			name:     "message with multiple spaces",
			content:  "   multiple   spaces   ",
			expected: []string{"multiple", "spaces"},
		},
		{
			name:     "UTF-8 message",
			content:  "Привет мир",
			expected: []string{"привет", "мир"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(tt.content)

			if msg.Content != tt.content {
				t.Errorf("Expected content %q, got %q", tt.content, msg.Content)
			}

			if !reflect.DeepEqual(msg.Words, tt.expected) {
				t.Errorf("Expected words %v, got %v", tt.expected, msg.Words)
			}
		})
	}
}
