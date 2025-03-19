package services

import (
	"context"
	"strconv"
	"time"

	"koneksi/services/iam/config"
	"koneksi/services/iam/core/logger"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
	prefix string // Add prefix field
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
		prefix: redisConfig.RedisPrefix, // Set prefix
	}
}

// GetRedis retrieves the Redis client instance
func (r *RedisService) GetRedis() *redis.Client {
	if r.client == nil {
		logger.Log.Fatal("redis not initialized")
	}
	return r.client
}

// AddPrefix adds the global prefix to a key
func (r *RedisService) AddPrefix(key string) string {
	return r.prefix + ":" + key
}
