package usecase

import (
	"aggregator/internal/domain"
	"aggregator/internal/repository"
	"log"
	"strings"
)

type StatsUseCase interface {
	ProcessMessage(content string)
	GetStats() (int, int)
}

type statsUseCase struct {
	statsRepo repository.StatsRepository
}

func NewStatsUseCase(statsRepo repository.StatsRepository) StatsUseCase {
	return &statsUseCase{
		statsRepo: statsRepo,
	}
}

func (uc *statsUseCase) ProcessMessage(content string) {
	if strings.TrimSpace(content) == "" {
		return
	}

	message := domain.NewMessage(content)
	uc.statsRepo.AddMessage(message)

	log.Printf("Processed message: %s", content)
}

func (uc *statsUseCase) GetStats() (int, int) {
	totalLines, uniqueWords := uc.statsRepo.GetStats()
	log.Printf("Stats requested: total_lines=%d, unique_words=%d", totalLines, uniqueWords)
	return totalLines, uniqueWords
}
