package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"insightiq/backend/internal/agent"
	"insightiq/backend/internal/connectors"
	httpserver "insightiq/backend/internal/http" // Fixed: Use alias to avoid conflict
	"insightiq/backend/internal/services"
)

func main() {
	// Health check flag for distroless image
	healthCheck := flag.Bool("health-check", false, "Perform health check and exit")
	flag.Parse()

	if *healthCheck {
		performHealthCheck()
		return
	}

	// Setup logging
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") == "true" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "password" || a.Key == "token" || a.Key == "secret" {
				return slog.String(a.Key, "***REDACTED***")
			}
			return a
		},
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize connectors
	supersetConn := connectors.NewSuperSetConnector(
		getEnvOrDefault("SUPERSET_URL", "http://superset:8088"),
		getEnvOrDefault("SUPERSET_USERNAME", "admin"),
		getEnvOrDefault("SUPERSET_PASSWORD", "admin"), logger)

	// Initialize PostgreSQL connector
	postgresConn, err := connectors.NewPostgresConnector(
		getEnvOrDefault("POSTGRES_URL", "postgres://superset:superset@postgres:5432/superset?sslmode=disable"),
		logger)
	if err != nil {
		logger.Error("Failed to initialize PostgreSQL connector", "error", err)
		os.Exit(1)
	}

	ollamaConn := connectors.NewOllamaConnector(
		getEnvOrDefault("OLLAMA_URL", "http://ollama:11434"), logger)

	whisperConn := connectors.NewWhisperConnector(
		getEnvOrDefault("WHISPER_URL", "http://whisper:9000"), logger)

	// Initialize agent manager
	agentManager := agent.NewManager(logger)

	// Create and register agents
	analyticsAgent := agent.NewAnalyticsAgent("analytics-1", supersetConn, postgresConn, ollamaConn, logger)
	voiceAgent := agent.NewVoiceAgent("voice-1", whisperConn, analyticsAgent, logger) // Fixed: This should work now

	if err := agentManager.RegisterAgent(analyticsAgent); err != nil {
		logger.Error("Failed to register analytics agent", "error", err)
		os.Exit(1)
	}

	if err := agentManager.RegisterAgent(voiceAgent); err != nil {
		logger.Error("Failed to register voice agent", "error", err)
		os.Exit(1)
	}

	// Start agent manager
	go func() {
		if err := agentManager.Start(ctx); err != nil {
			logger.Error("Agent manager failed to start", "error", err)
			cancel()
		}
	}()

	// Initialize services
	analyticsService := services.NewAnalyticsService(agentManager, logger)
	voiceService := services.NewVoiceService(agentManager, logger)

	// Create HTTP server
	httpServer := httpserver.NewServer(analyticsService, voiceService, logger) // Fixed: Use alias

	server := &http.Server{
		Addr:              getEnvOrDefault("PORT", ":8080"),
		Handler:           httpServer,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	// Start server
	go func() {
		logger.Info("Starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("Received shutdown signal")
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	// Close database connection
	if err := postgresConn.Close(); err != nil {
		logger.Error("Error closing PostgreSQL connection", "error", err)
	}

	logger.Info("Server stopped gracefully")
}

func performHealthCheck() {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8080/api/health")
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
	os.Exit(0)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
