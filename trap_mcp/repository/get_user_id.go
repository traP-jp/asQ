package repository

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	userCacheStore CacheStore[map[string]string]
	userInnerFn    func(context.Context) (map[string]string, error) = GetWithCache(
		updateCache,
		&userCacheStore,
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
		(*cache)[v.Name] = v.Id
	}
	return nil
}

func GetUserToId(ctx context.Context) (map[string]string, error) {
	return userInnerFn(ctx)
}
