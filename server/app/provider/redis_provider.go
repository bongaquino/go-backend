package provider

import (
	"context"
	"strconv"
	"time"

	"koneksi/server/config"
	"koneksi/server/core/logger"

	"github.com/go-redis/redis/v8"
)

type RedisProvider struct {
	client *redis.Client
	prefix string
}

// NewRedisProvider initializes a new RedisProvider
func NewRedisProvider() *RedisProvider {
	redisConfig := config.LoadRedisConfig()

	options := &redis.Options{
		Addr:     redisConfig.RedisHost + ":" + strconv.Itoa(redisConfig.RedisPort),
		Password: redisConfig.RedisPassword,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := redis.NewClient(options)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Log.Fatal("redis connection error", logger.Error(err))
	}

	return &RedisProvider{
		client: client,
		prefix: redisConfig.RedisPrefix,
	}
}

// prefixedKey adds the global prefix to a key if a prefix is set
func (r *RedisProvider) prefixedKey(key string) string {
	if r.prefix != "" {
		return r.prefix + ":" + key
	}
	return key
}

// Set sets a key-value pair in Redis with the given expiration
func (r *RedisProvider) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	prefixedKey := r.prefixedKey(key)
	return r.client.Set(ctx, prefixedKey, value, expiration).Err()
}

// Get retrieves the value of a key from Redis
func (r *RedisProvider) Get(ctx context.Context, key string) (string, error) {
	prefixedKey := r.prefixedKey(key)
	return r.client.Get(ctx, prefixedKey).Result()
}

// Del deletes a key from Redis
func (r *RedisProvider) Del(ctx context.Context, key string) error {
	prefixedKey := r.prefixedKey(key)
	return r.client.Del(ctx, prefixedKey).Err()
}
