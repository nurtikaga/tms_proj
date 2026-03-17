package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

const keyPrefix = "shipment:"

type RedisCache struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func key(id string) string { return keyPrefix + id }

func (c *RedisCache) Get(ctx context.Context, id string) (*entity.Shipment, error) {
	raw, err := c.client.Get(ctx, key(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis: get %s: %w", id, err)
	}
	var s entity.Shipment
	if err = json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("redis: unmarshal shipment: %w", err)
	}
	return &s, nil
}

func (c *RedisCache) Set(ctx context.Context, s *entity.Shipment, ttl time.Duration) error {
	raw, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("redis: marshal shipment: %w", err)
	}
	if err = c.client.Set(ctx, key(s.ID), raw, ttl).Err(); err != nil {
		return fmt.Errorf("redis: set %s: %w", s.ID, err)
	}
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, id string) error {
	if err := c.client.Del(ctx, key(id)).Err(); err != nil {
		return fmt.Errorf("redis: delete %s: %w", id, err)
	}
	return nil
}
