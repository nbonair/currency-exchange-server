package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nbonair/currency-exchange-server/configs"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataInSet(ctx context.Context, key string, data any) (bool, error)
}

func NewClient(cfg configs.CacheConfig) (Client, error) {
	switch cfg.Type {
	case configs.CacheTypeInMemory:
		return NewInMemoryClient(), nil
	case configs.CacheTypeRedis:
		return NewRedisClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cfg.Type)
	}
}

type redisClient struct {
	redisClient *redis.Client
}

func NewRedisClient(cfg configs.CacheConfig) Client {
	return &redisClient{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     cfg.Address,
			Username: cfg.Username,
			Password: cfg.Password,
		})}
}

func (c *redisClient) Get(ctx context.Context, key string) (any, error) {
	data, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get data from cache: %w", err) //status.Error(codes.Internal, "failed to get data into cache")
	}

	return data, nil
}

func (c *redisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if err := c.redisClient.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set data to cache: %w", err) //status.Error(codes.Internal, "failed to set data into cache")
	}
	return nil
}

func (c *redisClient) AddToSet(ctx context.Context, key string, data ...any) error {
	if err := c.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		return fmt.Errorf("failed to set data to set inside cache: %w", err)
	}
	return nil
}

// IsDataInSet implements Client.
func (c *redisClient) IsDataInSet(ctx context.Context, key string, data any) (bool, error) {
	res, err := c.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check data is member of set inside cache: %w", err)
	}
	return res, nil
}

type inMemoryClient struct {
	cache      map[string]any
	cacheMutex *sync.Mutex
}

func NewInMemoryClient() Client {
	return &inMemoryClient{
		cache:      make(map[string]any),
		cacheMutex: new(sync.Mutex),
	}
}

func (c *inMemoryClient) Get(ctx context.Context, key string) (any, error) {
	data, ok := c.cache[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	return data, nil
}

func (c *inMemoryClient) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	c.cache[key] = data
	return nil
}

func (c *inMemoryClient) AddToSet(_ context.Context, key string, data ...any) error {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	set := c.getSet(key)
	set = append(set, data...)
	c.cache[key] = set
	return nil
}

func (c *inMemoryClient) IsDataInSet(_ context.Context, key string, data any) (bool, error) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	set := c.getSet(key)

	for i := range set {
		if set[i] == data {
			return true, nil
		}
	}

	return false, nil
}

func (c *inMemoryClient) getSet(key string) []any {
	setValue, ok := c.cache[key]
	if !ok {
		return make([]any, 0)
	}

	set, ok := setValue.([]any)
	if !ok {
		return make([]any, 0)
	}

	return set
}
