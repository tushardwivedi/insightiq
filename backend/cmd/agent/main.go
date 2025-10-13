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

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"insightiq/backend/internal/agent"
	"insightiq/backend/internal/connectors"
	"insightiq/backend/internal/embedding"
	httpserver "insightiq/backend/internal/http" // Fixed: Use alias to avoid conflict
	"insightiq/backend/internal/intent"
	"insightiq/backend/internal/repository"
	"insightiq/backend/internal/schema"
	"insightiq/backend/internal/services"
	"insightiq/backend/internal/vectorstore"
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

	// Initialize PostgreSQL connector for sample data
	postgresURL := getEnvOrDefault("POSTGRES_URL", "postgres://insightiq_user:insightiq_password@postgres:5432/insightiq?sslmode=disable")
	postgresConn, err := connectors.NewPostgresConnector(postgresURL, logger)
	if err != nil {
		logger.Error("Failed to initialize PostgreSQL connector", "error", err)
		os.Exit(1)
	}

	// Initialize database connection for connector management
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		logger.Error("Failed to connect to database for connector management", "error", err)
		os.Exit(1)
	}

	// Initialize connector repository and service
	connectorRepo := repository.NewConnectorRepository(db)

	// Create connector tables if they don't exist
	if err := connectorRepo.CreateTables(ctx); err != nil {
		logger.Error("Failed to create connector tables", "error", err)
		os.Exit(1)
	}

	connectorService := services.NewConnectorService(connectorRepo, logger)

	// Initialize RAG infrastructure
	qdrantURL := getEnvOrDefault("QDRANT_URL", "http://qdrant:6333")
	vectorStore := vectorstore.NewQdrantClient(qdrantURL, logger)

	// Initialize embedding service
	ollamaURL := getEnvOrDefault("OLLAMA_URL", "http://ollama:11434")
	embeddingService := embedding.NewOllamaEmbeddingService(ollamaURL, "nomic-embed-text", logger)

	// Initialize schema ingestion service
	ingestionService := schema.NewIngestionService(vectorStore, embeddingService, logger)

	// Initialize intent classification service
	intentService := intent.NewClassificationService(vectorStore, embeddingService, logger)

	// Ingest domain contexts on startup
	go func() {
		if err := ingestionService.IngestDomainContexts(ctx); err != nil {
			logger.Error("Failed to ingest domain contexts", "error", err)
		} else {
			logger.Info("Domain contexts ingested successfully")
		}
	}()

	ollamaConn := connectors.NewOllamaConnector(ollamaURL, logger)

	whisperConn := connectors.NewWhisperConnector(
		getEnvOrDefault("WHISPER_URL", "http://whisper:9000"), logger)

	// Initialize agent manager
	agentManager := agent.NewManager(logger)

	// Create enhanced analytics service (connector-only architecture)
	enhancedAnalyticsService := services.NewEnhancedAnalyticsService(connectorService, ollamaConn, nil, nil, logger)

	// Create and register agents (PostgreSQL connections disabled - using connector-only architecture)
	analyticsAgent := agent.NewAnalyticsAgent("analytics-1", nil, nil, ollamaConn, logger)
	voiceAgent := agent.NewVoiceAgent("voice-1", whisperConn, analyticsAgent, logger)

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

	// Initialize services with enhanced analytics and RAG intent classification
	analyticsService := services.NewAnalyticsServiceWithRAG(agentManager, enhancedAnalyticsService, connectorService, ollamaConn, intentService, logger)
	voiceService := services.NewVoiceService(agentManager, logger)

	// Create planner service
	plannerService := services.NewPlannerService(ollamaConn, connectorService, logger)

	// Create HTTP server
	httpServer := httpserver.NewServer(analyticsService, voiceService, connectorService, plannerService, logger) // Fixed: Use alias

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

	// Close database connections
	if err := postgresConn.Close(); err != nil {
		logger.Error("Error closing PostgreSQL connection", "error", err)
	}

	if err := db.Close(); err != nil {
		logger.Error("Error closing database connection", "error", err)
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
