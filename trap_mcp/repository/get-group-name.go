package repository

import (
	"context"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

var (
	IdToNameCache   map[string]string
	updatedGroupsAt time.Time = time.UnixMicro(0)
)

func UpdateGroupCache(ctx context.Context) error {
	traq_client := clients.GetTraqClient()
	groups := traq_client.GroupApi.GetUserGroups(ctx)
	res, _, err := groups.Execute()
	if err != nil {
		return err
	}
	id_to_name := make(map[string]string)
	for _, v := range res {
		id_to_name[v.Id] = v.Name
	}
	IdToNameCache = id_to_name
	return nil
}

func GetGroupToName(ctx context.Context) (map[string]string, error) {
	now := time.Now()
	if now.Sub(updatedGroupsAt) > time.Hour {
		UpdateGroupCache(ctx)
	}

	return IdToNameCache, nil
}
