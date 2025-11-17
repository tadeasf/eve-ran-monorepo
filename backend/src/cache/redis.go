package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tadeasf/eve-ran/src/utils"
)

var (
	// RedisClient is the global Redis client instance
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis initializes the Redis connection
func InitRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "", // No password set
		DB:           0,  // Use default DB
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
	}

	utils.InfoLogger.Printf("Successfully connected to Redis at %s", addr)
	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// Get retrieves a value from Redis
func Get(key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	}
	return val, err
}

// Set stores a value in Redis with TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	return RedisClient.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key from Redis
func Delete(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// FlushAll clears all Redis data (use with caution)
func FlushAll() error {
	return RedisClient.FlushAll(ctx).Err()
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return RedisClient
}
