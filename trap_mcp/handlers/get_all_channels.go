package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/repository/id_to_channel"
)

func GetAllChannelsTool() mcp.Tool {
	tool := mcp.NewTool("getAllChannels",
		mcp.WithDescription("Get all channels(general)"),
	)
	return tool
}

func GetAllChannelsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resStr := ""
	idToChannel, err := id_to_channel.GetIdToChannel(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	for _, v := range idToChannel {
		resStr += v + "\n"

	}
	return mcp.NewToolResultText(resStr), nil
}
