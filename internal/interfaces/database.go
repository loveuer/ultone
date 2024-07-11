package interfaces

import (
	"context"
	"time"
)

type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error)
	// Set value 会被序列化, 优先使用 MarshalBinary 方法, 没有则执行 json.Marshal
	Set(ctx context.Context, key string, value any) error
	// SetEx value 会被序列化, 优先使用 MarshalBinary 方法, 没有则执行 json.Marshal
	SetEx(ctx context.Context, key string, value any, duration time.Duration) error
	Del(ctx context.Context, keys ...string) error
}
