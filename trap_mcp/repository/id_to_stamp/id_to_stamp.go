package id_to_stamp

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/util"
)

var (
	store   util.CacheStore[map[string]string]
	innerFn func(context.Context) (map[string]string, error) = util.GetWithCache(
		updateStampCache,
		&store,
		time.Minute*10,
	)
)

func updateStampCache(ctx context.Context, cache *map[string]string) error {
	traqClient := clients.GetTraqClient()
	channels, _, err := traqClient.StampApi.GetStamps(ctx).Execute()
	if err != nil {
		return err
	}
	(*cache) = make(map[string]string)
	for _, channel := range channels {
		(*cache)[channel.Id] = channel.Name
	}
	return nil
}

func GetIdToStamp(ctx context.Context) (map[string]string, error) {
	return innerFn(ctx)
}
