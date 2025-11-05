package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(addr, password string, db int, prefix string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		prefix: prefix,
	}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	fullKey := r.prefix + key
	return r.client.Set(ctx, fullKey, data, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := r.prefix + key
	data, err := r.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.prefix + key
	return r.client.Del(ctx, fullKey).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.prefix + key
	count, err := r.client.Exists(ctx, fullKey).Result()
	return count > 0, err
}

func (r *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	fullKey := r.prefix + key
	return r.client.SetNX(ctx, fullKey, data, ttl).Result()
}

func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	fullKey := r.prefix + key
	return r.client.Incr(ctx, fullKey).Result()
}

func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	fullKey := r.prefix + key
	return r.client.Expire(ctx, fullKey, ttl).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) Client() *redis.Client {
	return r.client
}
