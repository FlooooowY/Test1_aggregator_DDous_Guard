package domain

import "sync"

type Stats struct {
	TotalLines  int
	UniqueWords int
	Words       map[string]int
	Mu          sync.RWMutex
}

func NewStats() *Stats {
	return &Stats{
		Words: make(map[string]int),
	}
}

func (s *Stats) AddWord(word string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if word == "" {
		return
	}

	if _, exists := s.Words[word]; !exists {
		s.UniqueWords++
	}
	s.Words[word]++
}

func (s *Stats) AddLine() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.TotalLines++
}

func (s *Stats) GetStats() (int, int) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.TotalLines, s.UniqueWords
}
