package id_to_user

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/util"
)

var (
	store   util.CacheStore[map[string]string]
	innerFn func(context.Context) (map[string]string, error) = util.GetWithCache(
		updateCache,
		&store,
		time.Hour,
	)
)

func updateCache(ctx context.Context, cache *map[string]string) error {
	traq_client := clients.GetTraqClient()
	users, _, err := traq_client.UserApi.GetUsers(ctx).IncludeSuspended(true).Execute()
	if err != nil {
		return err
	}
	*cache = make(map[string]string)
	for _, v := range users {
		(*cache)[v.Id] = v.Name
	}
	return nil
}

func GetIdToUserId(ctx context.Context) (map[string]string, error) {
	return innerFn(ctx)
}
