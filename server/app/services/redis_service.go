package services

import (
	"context"
	"strconv"
	"time"

	"koneksi/server/config"
	"koneksi/server/core/logger"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
	prefix string
}

// NewRedisService initializes a new RedisService
func NewRedisService() *RedisService {
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

	return &RedisService{
		client: client,
		prefix: redisConfig.RedisPrefix,
	}
}

// prefixedKey adds the global prefix to a key if a prefix is set
func (r *RedisService) prefixedKey(key string) string {
	if r.prefix != "" {
		return r.prefix + ":" + key
	}
	return key
}

// Set sets a key-value pair in Redis with the given expiration
func (r *RedisService) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	prefixedKey := r.prefixedKey(key)
	return r.client.Set(ctx, prefixedKey, value, expiration).Err()
}

// Get retrieves the value of a key from Redis
func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	prefixedKey := r.prefixedKey(key)
	return r.client.Get(ctx, prefixedKey).Result()
}

// Del deletes a key from Redis
func (r *RedisService) Del(ctx context.Context, key string) error {
	prefixedKey := r.prefixedKey(key)
	return r.client.Del(ctx, prefixedKey).Err()
}
