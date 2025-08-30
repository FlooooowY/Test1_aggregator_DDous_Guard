package repository

import (
	"aggregator/internal/domain"
	"sync"
)

type StatsRepository interface {
	AddMessage(message *domain.Message)
	GetStats() (int, int)
}

type inMemoryStatsRepository struct {
	stats *domain.Stats
	mu    sync.RWMutex
}

func NewInMemoryStatsRepository() StatsRepository {
	return &inMemoryStatsRepository{
		stats: domain.NewStats(),
	}
}

func (r *inMemoryStatsRepository) AddMessage(message *domain.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stats.AddLine()

	for _, word := range message.Words {
		r.stats.AddWord(word)
	}
}

func (r *inMemoryStatsRepository) GetStats() (int, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.stats.GetStats()
}
