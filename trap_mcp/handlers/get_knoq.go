package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repository/channel_to_id"
	"github.com/traP-jp/h25s_05/trap_mcp/repository/user_to_id"
)

func GetKnoqTool() mcp.Tool {
	tool := mcp.NewTool("get_progress_room",
		mcp.WithDescription("Get progress room"),
	)
	return tool
}

func GetKnoqHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	channelIdMap, err := channel_to_id.GetChannelToId(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	userIdMap, err := user_to_id.GetUserToId(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	traq_client := clients.GetTraqClient()
	getReq := traq_client.ChannelApi.GetMessages(ctx, channelIdMap["services/knoQ/daily"])
	fmt.Print(channelIdMap["services/knoQ/daily"])

	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	i := 0
	for ; ; i++ {
		if res[i].UserId == userIdMap["Webhook#qzrMONHOQSShr_OjaiEMhQ"] {
			break
		}
	}

	return mcp.NewToolResultText(res[i].Content), nil
}
