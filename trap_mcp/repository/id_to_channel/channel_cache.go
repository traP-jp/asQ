package id_to_channel

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/util"
)

var (
	store   util.CacheStore[map[string]string]
	innerFn func(context.Context) (map[string]string, error) = util.GetWithCache(
		updateChannelCache,
		&store,
		time.Minute*10,
	)
)

func updateChannelCache(ctx context.Context, cache *map[string]string) error {
	traqClient := clients.GetTraqClient()
	channels, _, err := traqClient.ChannelApi.GetChannels(ctx).Execute()
	if err != nil {
		return err
	}
	channelsWithPath, err := util.ConvertChannelNameToPath(channels)
	if err != nil {
		return err
	}

	*cache = make(map[string]string)
	for _, channel := range channelsWithPath.Public {
		(*cache)[channel.Id] = channel.Name
	}
	return nil
}

func GetIdToChannel(ctx context.Context) (map[string]string, error) {
	return innerFn(ctx)
}
