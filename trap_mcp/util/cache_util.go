package util

import (
	"context"
	"sync"
	"time"
)

type CacheStore[T any] struct {
	Data      T
	UpdatedAt time.Time
	Mutex     sync.Mutex
}

func GetWithCache[T any](
	update_fn func(context.Context, *T) error,
	store *CacheStore[T],
	interval time.Duration,
) func(ctx context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		store.Mutex.Lock()
		defer store.Mutex.Unlock()

		if time.Since(store.UpdatedAt) > interval {
			if err := update_fn(ctx, &store.Data); err != nil {
				return store.Data, err
			}
			store.UpdatedAt = time.Now()
		}

		return store.Data, nil
	}
}
