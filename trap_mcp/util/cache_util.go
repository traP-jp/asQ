package util

import (
	"context"
	"sync"
	"time"
)

type CacheStore[T any] struct {
	Data      T
	UpdatedAt time.Time
	RwLock    sync.RWMutex
}

func GetWithCache[T any](
	update_fn func(context.Context, *T) error,
	store *CacheStore[T],
	interval time.Duration,
) func(ctx context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		var err error = nil
		if time.Since(store.UpdatedAt) > interval {
			store.RwLock.Lock()
			err = update_fn(ctx, &store.Data)
			store.RwLock.Unlock()
			store.UpdatedAt = time.Now()
		}
		store.RwLock.RLock()
		data := store.Data
		store.RwLock.RUnlock()
		return data, err
	}
}
