package cache

import (
	"context"
	"fmt"
	"time"
	"ultone/internal/interfaces"

	"gitea.com/taozitaozi/gredis"
)

var _ interfaces.Cacher = (*_mem)(nil)

type _mem struct {
	client *gredis.Gredis
}

func (m *_mem) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}

	bs, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid value type=%T", v)
	}

	return bs, nil
}

func (m *_mem) GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error) {
	v, err := m.client.GetEx(key, duration)
	if err != nil {
		return nil, err
	}

	bs, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid value type=%T", v)
	}

	return bs, nil
}

func (m *_mem) Set(ctx context.Context, key string, value any) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}
	return m.client.Set(key, bs)
}

func (m *_mem) SetEx(ctx context.Context, key string, value any, duration time.Duration) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}
	return m.client.SetEx(key, bs, duration)
}

func (m *_mem) Del(ctx context.Context, keys ...string) error {
	m.client.Delete(keys...)
	return nil
}
