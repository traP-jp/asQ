package repository

import (
	"context"
	"sync"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	idToStampName      map[string]string
	idToStampUpdatedAt time.Time = time.UnixMicro(0)
	stampCacheMutex    sync.Mutex
)

func updateStampCache(ctx context.Context) error {
	stampCacheMutex.Lock()
	defer stampCacheMutex.Unlock()
	if time.Since(idToStampUpdatedAt) < 10*time.Minute {
		return nil
	}
	traqClient := clients.GetTraqClient()
	channels, _, err := traqClient.StampApi.GetStamps(ctx).Execute()
	if err != nil {
		return err
	}
	idToStampName = make(map[string]string)
	for _, channel := range channels {
		idToStampName[channel.Id] = channel.Name
	}
	idToStampUpdatedAt = time.Now()
	return nil
}

func GetIdToStamp(ctx context.Context) (map[string]string, error) {
	if err := updateStampCache(ctx); err != nil {
		return nil, err
	}
	return idToStampName, nil
}
