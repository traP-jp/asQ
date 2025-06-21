package repository

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	groupCacheStore CacheStore[map[string]string]
	groupInnerFn    func(context.Context) (map[string]string, error) = GetWithCache(
		UpdateGroupCache,
		&groupCacheStore,
		time.Hour,
	)
)

func UpdateGroupCache(ctx context.Context, cache *map[string]string) error {
	traq_client := clients.GetTraqClient()
	groups := traq_client.GroupApi.GetUserGroups(ctx)
	res, _, err := groups.Execute()
	if err != nil {
		return err
	}
	*cache = make(map[string]string)
	for _, v := range res {
		(*cache)[v.Id] = v.Name
	}
	return nil
}

func GetGroupToName(ctx context.Context) (map[string]string, error) {
	return groupInnerFn(ctx)
}
