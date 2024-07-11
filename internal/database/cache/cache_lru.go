package cache

import (
	"context"
	"github.com/hashicorp/golang-lru/v2/expirable"
	_ "github.com/hashicorp/golang-lru/v2/expirable"
	"time"
	"ultone/internal/interfaces"
)

var _ interfaces.Cacher = (*_lru)(nil)

type _lru struct {
	client *expirable.LRU[string, *_lru_value]
}

type _lru_value struct {
	duration time.Duration
	last     time.Time
	bs       []byte
}

func (l *_lru) Get(ctx context.Context, key string) ([]byte, error) {
	v, ok := l.client.Get(key)
	if !ok {
		return nil, ErrorKeyNotFound
	}

	if v.duration == 0 {
		return v.bs, nil
	}

	if time.Now().Sub(v.last) > v.duration {
		l.client.Remove(key)
		return nil, ErrorKeyNotFound
	}

	return v.bs, nil
}

func (l *_lru) GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error) {
	v, ok := l.client.Get(key)
	if !ok {
		return nil, ErrorKeyNotFound
	}

	if v.duration == 0 {
		return v.bs, nil
	}

	now := time.Now()

	if now.Sub(v.last) > v.duration {
		l.client.Remove(key)
		return nil, ErrorKeyNotFound
	}

	l.client.Add(key, &_lru_value{
		duration: duration,
		last:     now,
		bs:       v.bs,
	})

	return v.bs, nil
}

func (l *_lru) Set(ctx context.Context, key string, value any) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}

	l.client.Add(key, &_lru_value{
		duration: 0,
		last:     time.Now(),
		bs:       bs,
	})

	return nil
}

func (l *_lru) SetEx(ctx context.Context, key string, value any, duration time.Duration) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}

	l.client.Add(key, &_lru_value{
		duration: duration,
		last:     time.Now(),
		bs:       bs,
	})

	return nil
}

func (l *_lru) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		l.client.Remove(key)
	}

	return nil
}

func newLRUCache() (interfaces.Cacher, error) {
	client := expirable.NewLRU[string, *_lru_value](0, nil, 0)

	return &_lru{client: client}, nil
}
