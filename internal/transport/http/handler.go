package http

import (
	"aggregator/internal/usecase"
	"encoding/json"
	"net/http"
)

type StatsResponse struct {
	TotalLines  int `json:"total_lines"`
	UniqueWords int `json:"unique_words"`
}

type Handler struct {
	statsUseCase usecase.StatsUseCase
}

func NewHandler(statsUseCase usecase.StatsUseCase) *Handler {
	return &Handler{
		statsUseCase: statsUseCase,
	}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	totalLines, uniqueWords := h.statsUseCase.GetStats()

	response := StatsResponse{
		TotalLines:  totalLines,
		UniqueWords: uniqueWords,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
