package repository

import (
	"context"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func GetUserToId(ctx context.Context) (map[string]string, error) {
	traq_client := clients.GetTraqClient()
	users, _, err := traq_client.UserApi.GetUsers(ctx).IncludeSuspended(true).Execute()
	if err != nil {
		return nil, err
	}
	name_to_id := make(map[string]string)
	for _, v := range users {
		name_to_id[v.Name] = v.Id
	}

	return name_to_id, nil
}
