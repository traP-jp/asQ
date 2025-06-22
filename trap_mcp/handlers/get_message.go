package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func GetMessageTool() mcp.Tool {
	tool := mcp.NewTool("getMessage",
		mcp.WithDescription("Get message from ID"),
		mcp.WithString("messageId",
			mcp.Description("Message ID to get"),
			mcp.DefaultString(""),
		),
	)
	return tool
}

func GetMessageHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	messageId, err := request.RequireString("messageId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	getReq := traq_client.MessageApi.GetMessage(ctx, messageId)
	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res.Content), nil
}
