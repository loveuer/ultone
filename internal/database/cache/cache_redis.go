package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type _redis struct {
	client *redis.Client
}

func (r *_redis) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func (r *_redis) GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error) {
	result, err := r.client.GetEx(ctx, key, duration).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func (r *_redis) Set(ctx context.Context, key string, value any) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}

	_, err = r.client.Set(ctx, key, bs, redis.KeepTTL).Result()
	return err
}

func (r *_redis) SetEx(ctx context.Context, key string, value any, duration time.Duration) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}

	_, err = r.client.SetEX(ctx, key, bs, duration).Result()

	return err
}

func (r *_redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
