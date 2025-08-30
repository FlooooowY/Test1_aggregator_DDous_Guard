package main

import (
	"aggregator/internal/app"
	"aggregator/internal/config"
	"log"
	"os"
)

func main() {
	cfg := config.Load()

	application := app.NewApp(cfg)

	if err := application.Run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}

	log.Println("Application stopped gracefully")
}
