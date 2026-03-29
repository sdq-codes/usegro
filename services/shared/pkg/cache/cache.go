package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache wraps a Redis client with a key prefix.
type Cache struct {
	rdb    redis.Cmdable
	prefix string
}

// New creates a Cache with the given Redis client and key prefix.
func New(rdb redis.Cmdable, prefix string) *Cache {
	return &Cache{rdb: rdb, prefix: prefix}
}

func (c *Cache) addPrefix(key string) string {
	return fmt.Sprintf("%s_%s", c.prefix, key)
}

// Set sets a key-value pair with an expiration time.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	key = c.addPrefix(key)
	if err := c.rdb.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Get retrieves the value of a key.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	key = c.addPrefix(key)
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// Pull retrieves the value of a key and then deletes it.
func (c *Cache) Pull(ctx context.Context, key string) (string, error) {
	key = c.addPrefix(key)
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	if _, delErr := c.rdb.Del(ctx, key).Result(); delErr != nil {
		return "", fmt.Errorf("failed to delete key %s: %w", key, delErr)
	}
	return val, nil
}

// SetForever sets a key-value pair with no expiration.
func (c *Cache) SetForever(ctx context.Context, key string, value interface{}) error {
	key = c.addPrefix(key)
	if err := c.rdb.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("failed to set key %s forever: %w", key, err)
	}
	return nil
}

// Remove deletes the key-value pair.
func (c *Cache) Remove(ctx context.Context, key string) error {
	key = c.addPrefix(key)
	if _, err := c.rdb.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to forget key %s: %w", key, err)
	}
	return nil
}

// Flush removes all keys matching the prefix pattern.
func (c *Cache) Flush(ctx context.Context) error {
	key := c.addPrefix("*")
	if _, err := c.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}

// Increment increases the integer value of a key by the given amount.
func (c *Cache) Increment(ctx context.Context, key string, increment int64) (int64, error) {
	key = c.addPrefix(key)
	val, err := c.rdb.IncrBy(ctx, key, increment).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s by %d: %w", key, increment, err)
	}
	return val, nil
}

// Decrement decreases the integer value of a key by the given amount.
func (c *Cache) Decrement(ctx context.Context, key string, decrement int64) (int64, error) {
	key = c.addPrefix(key)
	val, err := c.rdb.DecrBy(ctx, key, decrement).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s by %d: %w", key, decrement, err)
	}
	return val, nil
}

// Remember returns the cached value if it exists; otherwise calls fetchFunc, stores the result, and returns it.
func (c *Cache) Remember(ctx context.Context, key string, duration time.Duration, fetchFunc func() ([]byte, error)) ([]byte, error) {
	value, err := c.Get(ctx, key)
	if err == nil {
		return []byte(value), nil
	}

	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	if err := c.Set(ctx, key, data, duration); err != nil {
		return nil, err
	}
	return data, nil
}

// RememberForever is like Remember but stores the result without expiry.
func (c *Cache) RememberForever(ctx context.Context, key string, fetchFunc func() ([]byte, error)) ([]byte, error) {
	return c.Remember(ctx, key, 0, fetchFunc)
}
