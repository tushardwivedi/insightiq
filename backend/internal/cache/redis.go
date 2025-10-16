package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache provides caching functionality using Redis
type RedisCache struct {
	client *redis.Client
	logger *slog.Logger
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(redisURL string, logger *slog.Logger) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("✅ Redis cache connected successfully", "url", redisURL)

	return &RedisCache{
		client: client,
		logger: logger.With("component", "redis_cache"),
	}, nil
}

// Get retrieves a value from cache
func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := rc.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss, not an error
	}
	if err != nil {
		rc.logger.Error("Failed to get from cache", "key", key, "error", err)
		return nil, err
	}

	rc.logger.Debug("Cache hit", "key", key)
	return val, nil
}

// Set stores a value in cache with TTL
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data []byte
	var err error

	// Handle different types
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		// Marshal to JSON for complex types
		data, err = json.Marshal(value)
		if err != nil {
			rc.logger.Error("Failed to marshal value", "key", key, "error", err)
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	err = rc.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		rc.logger.Error("Failed to set cache", "key", key, "error", err)
		return err
	}

	rc.logger.Debug("Cache set", "key", key, "ttl", ttl)
	return nil
}

// Delete removes a value from cache
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		rc.logger.Error("Failed to delete from cache", "key", key, "error", err)
		return err
	}

	rc.logger.Debug("Cache deleted", "key", key)
	return nil
}

// DeletePattern removes all keys matching a pattern
func (rc *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	deletedCount := 0
	for iter.Next(ctx) {
		err := rc.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			rc.logger.Warn("Failed to delete key", "key", iter.Val(), "error", err)
			continue
		}
		deletedCount++
	}

	if err := iter.Err(); err != nil {
		rc.logger.Error("Failed to scan keys", "pattern", pattern, "error", err)
		return err
	}

	rc.logger.Info("Cache pattern deleted", "pattern", pattern, "count", deletedCount)
	return nil
}

// GetJSON retrieves and unmarshals a JSON value from cache
func (rc *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := rc.Get(ctx, key)
	if err != nil {
		return err
	}

	if data == nil {
		return nil // Cache miss
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		rc.logger.Error("Failed to unmarshal cached value", "key", key, "error", err)
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// Ping checks if Redis is alive
func (rc *RedisCache) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}

// Stats returns cache statistics
func (rc *RedisCache) Stats(ctx context.Context) (map[string]interface{}, error) {
	info, err := rc.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	// Get memory stats
	memInfo, err := rc.client.Info(ctx, "memory").Result()
	if err != nil {
		return nil, err
	}

	// Get keyspace stats
	keyspaceInfo, err := rc.client.Info(ctx, "keyspace").Result()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stats":    info,
		"memory":   memInfo,
		"keyspace": keyspaceInfo,
	}, nil
}
