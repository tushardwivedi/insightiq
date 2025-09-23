// cmd/agent/main.go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"insightiq/backend/internal/agent"
	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/http"
	"insightiq/backend/internal/services"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize connectors.
	supersetConn := connectors.NewSuperSetConnector(
		os.Getenv("SUPERSET_URL"), "admin", "admin", logger)

	ollamaConn := connectors.NewOllamaConnector(
		os.Getenv("OLLAMA_URL"), logger)

	whisperConn := connectors.NewWhisperConnector(
		os.Getenv("WHISPER_URL"), logger)

	// Initialize agent manager
	agentManager := agent.NewManager(logger)

	// Create and register agents
	analyticsAgent := agent.NewAnalyticsAgent("analytics-1", supersetConn, ollamaConn, logger)
	voiceAgent := agent.NewVoiceAgent("voice-1", whisperConn, analyticsAgent, logger)

	agentManager.RegisterAgent(analyticsAgent)
	agentManager.RegisterAgent(voiceAgent)

	// Start agent manager
	go agentManager.Start(ctx)

	// Initialize services
	analyticsService := services.NewAnalyticsService(agentManager, logger)
	voiceService := services.NewVoiceService(agentManager, logger)

	// Create HTTP server
	httpServer := http.NewServer(analyticsService, voiceService, logger)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      httpServer,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server stopped")
}
