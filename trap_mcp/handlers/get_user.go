package handlers

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repository"
)

func GetUserTool() mcp.Tool {
	tool := mcp.NewTool("getUser",
		mcp.WithDescription("Get user"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("user name"),
		),
	)
	return tool
}

func GetUserHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	name_to_id, err := repository.GetUserId(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id, ok := name_to_id[name]
	if !ok {
		return mcp.NewToolResultText("Username not found"), nil
	}
	getReq := traq_client.UserApi.GetUser(ctx, id)

	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	groupNameMap, err := repository.GetGroupToName(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	for i := 0; i < len(res.Groups); i++ {
		res.Groups[i] = groupNameMap[res.Groups[i]]
	}
	homeChannelMap, err := repository.GetChannelName(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if res.HomeChannel.IsSet() {
		homeChannelId := res.GetHomeChannel()
		homeChannelName := homeChannelMap[homeChannelId]
		res.SetHomeChannel(homeChannelName)
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
