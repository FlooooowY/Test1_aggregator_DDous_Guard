package app

import (
	"aggregator/internal/config"
	"aggregator/internal/repository"
	"aggregator/internal/transport/http"
	"aggregator/internal/transport/tcp"
	"aggregator/internal/usecase"
	"context"
	"fmt"
	"log"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	tcpHandler  *tcp.Handler
	httpHandler *http.Handler
}

type App struct {
	config      *config.Config
	server      *Server
	quit        chan os.Signal
	httpServer  *nethttp.Server
	tcpListener net.Listener
}

func NewApp(cfg *config.Config) *App {
	statsRepo := repository.NewInMemoryStatsRepository()
	statsUseCase := usecase.NewStatsUseCase(statsRepo)

	server := &Server{
		tcpHandler:  tcp.NewHandler(statsUseCase),
		httpHandler: http.NewHandler(statsUseCase),
	}

	return &App{
		config: cfg,
		server: server,
		quit:   make(chan os.Signal, 1),
	}
}

func (a *App) Run() error {
	signal.Notify(a.quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.startTCPServer(); err != nil {
			errChan <- fmt.Errorf("TCP server error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.startHTTPServer(); err != nil && err != nethttp.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %v", err)
		}
	}()

	a.logStartupInfo()

	<-a.quit
	log.Println("Received shutdown signal")

	return a.shutdown()
}

func (a *App) startTCPServer() error {
	listener, err := net.Listen("tcp", ":"+a.config.TCP.Port)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %v", err)
	}
	a.tcpListener = listener

	log.Printf("TCP server listening on port %s", a.config.TCP.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-context.Background().Done():
				return nil
			default:
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
		}

		go a.server.tcpHandler.HandleConnection(conn)
	}
}

func (a *App) startHTTPServer() error {
	mux := nethttp.NewServeMux()
	mux.HandleFunc("/stats", a.server.httpHandler.GetStats)

	a.httpServer = &nethttp.Server{
		Addr:    ":" + a.config.HTTP.Port,
		Handler: mux,
	}

	log.Printf("HTTP server listening on port %s", a.config.HTTP.Port)
	return a.httpServer.ListenAndServe()
}

func (a *App) shutdown() error {
	log.Println("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if a.httpServer != nil {
			if err := a.httpServer.Shutdown(ctx); err != nil {
				errChan <- fmt.Errorf("HTTP server shutdown error: %v", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if a.tcpListener != nil {
			if err := a.tcpListener.Close(); err != nil {
				errChan <- fmt.Errorf("TCP server shutdown error: %v", err)
			}
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}

	log.Println("Application shut down successfully")
	return nil
}

func (a *App) logStartupInfo() {
	log.Println("Starting servers...")
	log.Printf("TCP server: localhost:%s", a.config.TCP.Port)
	log.Printf("HTTP server: localhost:%s", a.config.HTTP.Port)
	log.Printf("Try: curl http://localhost:%s/stats", a.config.HTTP.Port)
	log.Println("Press Ctrl+C to shutdown gracefully")
}
