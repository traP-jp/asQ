package repository

import (
	"context"
	"sync"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	nameToUserIdCache  map[string]string
	idToUserNameCache  map[string]string
	UserCacheUpdatedAt time.Time = time.UnixMicro(0)
	userCacheMutex     sync.Mutex
)

func updateCache(ctx context.Context) error {
	userCacheMutex.Lock()
	defer userCacheMutex.Unlock()
	traq_client := clients.GetTraqClient()
	users, _, err := traq_client.UserApi.GetUsers(ctx).IncludeSuspended(true).Execute()
	if err != nil {
		return err
	}
	name_to_id := make(map[string]string)
	id_to_name := make(map[string]string)
	for _, v := range users {
		name_to_id[v.Name] = v.Id
		id_to_name[v.Id] = v.Name
	}
	nameToUserIdCache = name_to_id
	idToUserNameCache = id_to_name
	return nil
}

func GetUserId(ctx context.Context) (map[string]string, error) {
	now := time.Now()
	if now.Sub(UserCacheUpdatedAt) > time.Hour {
		updateCache(ctx)
	}

	return nameToUserIdCache, nil
}

func GetUserName(ctx context.Context) (map[string]string, error) {
	now := time.Now()
	if now.Sub(UserCacheUpdatedAt) > time.Hour {
		updateCache(ctx)
	}

	return idToUserNameCache, nil
}
