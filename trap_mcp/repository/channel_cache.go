package repository

import (
	"context"
	"sync"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	idToChannelName map[string]string
	updatedAt       time.Time = time.UnixMicro(0)
	cacheMutex      sync.Mutex
)

func updateChannelCache(ctx context.Context) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Since(updatedAt) < 10*time.Minute {
		return nil
	}
	traqClient := clients.GetTraqClient()
	channels, _, err := traqClient.ChannelApi.GetChannels(ctx).Execute()
	if err != nil {
		return err
	}
	idToChannelName = make(map[string]string)
	for _, channel := range channels.Public {
		idToChannelName[channel.Id] = channel.Name
	}
	updatedAt = time.Now()
	return nil
}

func GetIdToChannel(ctx context.Context) (map[string]string, error) {
	if err := updateChannelCache(ctx); err != nil {
		return nil, err
	}
	return idToChannelName, nil
}
