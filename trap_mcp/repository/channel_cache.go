package repository

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	channelCacheStore CacheStore[map[string]string]
	channelInnerFn    func(context.Context) (map[string]string, error) = GetWithCache(
		updateChannelCache,
		&channelCacheStore,
		time.Minute*10,
	)
)

func updateChannelCache(ctx context.Context, cache *map[string]string) error {
	traqClient := clients.GetTraqClient()
	channels, _, err := traqClient.ChannelApi.GetChannels(ctx).Execute()
	if err != nil {
		return err
	}
	*cache = make(map[string]string)
	for _, channel := range channels.Public {
		(*cache)[channel.Id] = channel.Name
	}
	return nil
}

func GetIdToChannel(ctx context.Context) (map[string]string, error) {
	return channelInnerFn(ctx)
}
