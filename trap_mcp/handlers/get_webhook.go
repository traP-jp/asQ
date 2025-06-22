package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func GetWebhookTool() mcp.Tool {
	tool := mcp.NewTool("getWebhook",
		mcp.WithDescription("Get info from services/knoQ/daily"),
	)
	return tool
}

func GetWebhookHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	getReq := traq_client.UserApi.GetUsers(ctx)
	useSuspended, err := request.RequireBool("includeSuspended")
	if err == nil {
		getReq = getReq.IncludeSuspended(useSuspended)
	}
	//res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(""), nil
}
