package repositry

import (
	"context"
	"fmt"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	nameToIdCache map[string]string
	updatedAt     time.Time = time.UnixMicro(0)
)

func updateCache(ctx context.Context) error {
	traq_client := clients.GetTraqClient()
	users, _, err := traq_client.UserApi.GetUsers(ctx).IncludeSuspended(true).Execute()
	if err != nil {
		return err
	}
	name_to_id := make(map[string]string)
	for _, v := range users {
		name_to_id[v.Name] = v.Id
	}
	nameToIdCache = name_to_id
	return nil
}

func GetUserToId(ctx context.Context) (map[string]string, error) {
	now := time.Now()
	if now.Sub(updatedAt) > time.Hour {
		updateCache(ctx)
	}

	return nameToIdCache, nil
}
