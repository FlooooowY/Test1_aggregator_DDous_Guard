package domain

import (
	"regexp"
	"strings"
)

type Message struct {
	Content string
	Words   []string
}

func NewMessage(content string) *Message {
	words := extractWords(content)
	return &Message{
		Content: content,
		Words:   words,
	}
}

func extractWords(content string) []string {
	if content == "" {
		return []string{}
	}

	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	parts := re.Split(content, -1)

	var words []string
	for _, part := range parts {
		word := strings.TrimSpace(part)
		if word != "" {
			words = append(words, strings.ToLower(word))
		}
	}

	return words
}
