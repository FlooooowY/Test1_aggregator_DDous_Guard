package tcp

import (
	"aggregator/internal/usecase"
	"bufio"
	"log"
	"net"
)

type Handler struct {
	statsUseCase usecase.StatsUseCase
}

func NewHandler(statsUseCase usecase.StatsUseCase) *Handler {
	return &Handler{
		statsUseCase: statsUseCase,
	}
}

func (h *Handler) HandleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		h.statsUseCase.ProcessMessage(message)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}
