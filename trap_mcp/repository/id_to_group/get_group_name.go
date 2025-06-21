package id_to_group

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/util"
)

var (
	store   util.CacheStore[map[string]string]
	innerFn func(context.Context) (map[string]string, error) = util.GetWithCache(
		UpdateGroupCache,
		&store,
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
	return innerFn(ctx)
}
