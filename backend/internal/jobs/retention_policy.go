package jobs

import (
	"context"
	"log/slog"
	"time"

	"insightiq/backend/internal/usecases"
)

// RetentionPolicy defines the interface for data retention policies
type RetentionPolicy interface {
	Run(ctx context.Context) error
	Start(ctx context.Context) error
	Stop() error
}

// QueryHistoryRetentionPolicy handles automatic cleanup of old query history
// based on configurable retention periods
type QueryHistoryRetentionPolicy struct {
	useCase          *usecases.QueryHistoryUseCase
	logger           *slog.Logger
	retentionPeriod  time.Duration
	checkInterval    time.Duration
	ticker           *time.Ticker
	stopChan         chan struct{}
	isRunning        bool
}

// RetentionPolicyConfig holds configuration for the retention policy
type RetentionPolicyConfig struct {
	// RetentionPeriod is how long to keep query history (default: 90 days)
	RetentionPeriod time.Duration

	// CheckInterval is how often to run the cleanup job (default: 24 hours)
	CheckInterval time.Duration
}

// NewQueryHistoryRetentionPolicy creates a new retention policy job
func NewQueryHistoryRetentionPolicy(
	useCase *usecases.QueryHistoryUseCase,
	config RetentionPolicyConfig,
	logger *slog.Logger,
) *QueryHistoryRetentionPolicy {
	// Set defaults if not provided
	if config.RetentionPeriod == 0 {
		config.RetentionPeriod = 90 * 24 * time.Hour // 90 days default
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = 24 * time.Hour // Daily checks
	}

	return &QueryHistoryRetentionPolicy{
		useCase:         useCase,
		logger:          logger.With("component", "retention_policy"),
		retentionPeriod: config.RetentionPeriod,
		checkInterval:   config.CheckInterval,
		stopChan:        make(chan struct{}),
		isRunning:       false,
	}
}

// Run executes a single cleanup cycle
func (rp *QueryHistoryRetentionPolicy) Run(ctx context.Context) error {
	rp.logger.Info("Running query history retention policy",
		"retention_period", rp.retentionPeriod.String())

	startTime := time.Now()

	// Delete old queries
	deletedCount, err := rp.useCase.DeleteOldQueries(ctx, rp.retentionPeriod)
	if err != nil {
		rp.logger.Error("Retention policy execution failed",
			"error", err,
			"duration", time.Since(startTime))
		return err
	}

	rp.logger.Info("Retention policy completed successfully",
		"deleted_count", deletedCount,
		"retention_period", rp.retentionPeriod.String(),
		"duration", time.Since(startTime))

	return nil
}

// Start begins the scheduled retention policy job
func (rp *QueryHistoryRetentionPolicy) Start(ctx context.Context) error {
	if rp.isRunning {
		rp.logger.Warn("Retention policy is already running")
		return nil
	}

	rp.logger.Info("Starting query history retention policy",
		"retention_period", rp.retentionPeriod.String(),
		"check_interval", rp.checkInterval.String())

	rp.ticker = time.NewTicker(rp.checkInterval)
	rp.isRunning = true

	// Run immediately on start
	go func() {
		// First execution
		if err := rp.Run(ctx); err != nil {
			rp.logger.Error("Initial retention policy run failed", "error", err)
		}

		// Then run on schedule
		for {
			select {
			case <-rp.ticker.C:
				if err := rp.Run(ctx); err != nil {
					rp.logger.Error("Scheduled retention policy run failed", "error", err)
				}

			case <-rp.stopChan:
				rp.logger.Info("Retention policy stopped")
				return

			case <-ctx.Done():
				rp.logger.Info("Retention policy context cancelled")
				return
			}
		}
	}()

	return nil
}

// Stop halts the scheduled retention policy job
func (rp *QueryHistoryRetentionPolicy) Stop() error {
	if !rp.isRunning {
		return nil
	}

	rp.logger.Info("Stopping query history retention policy")

	close(rp.stopChan)

	if rp.ticker != nil {
		rp.ticker.Stop()
	}

	rp.isRunning = false

	rp.logger.Info("Retention policy stopped successfully")
	return nil
}

// IsRunning returns whether the retention policy is currently active
func (rp *QueryHistoryRetentionPolicy) IsRunning() bool {
	return rp.isRunning
}

// GetRetentionPeriod returns the current retention period
func (rp *QueryHistoryRetentionPolicy) GetRetentionPeriod() time.Duration {
	return rp.retentionPeriod
}

// GetCheckInterval returns the current check interval
func (rp *QueryHistoryRetentionPolicy) GetCheckInterval() time.Duration {
	return rp.checkInterval
}

// UpdateRetentionPeriod dynamically updates the retention period
// This does not require restarting the job
func (rp *QueryHistoryRetentionPolicy) UpdateRetentionPeriod(period time.Duration) {
	if period < 24*time.Hour {
		rp.logger.Warn("Retention period too short, ignoring update",
			"requested", period.String(),
			"minimum", "24h")
		return
	}

	oldPeriod := rp.retentionPeriod
	rp.retentionPeriod = period

	rp.logger.Info("Retention period updated",
		"old_period", oldPeriod.String(),
		"new_period", period.String())
}
