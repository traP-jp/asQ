package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func GetAllUsrsTool() mcp.Tool {
	tool := mcp.NewTool("getAllUsers",
		mcp.WithDescription("Get all traP users"),
	)
	return tool
}

func GetAllUsersHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	getReq := traq_client.UserApi.GetUsers(ctx)
	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	usersStr := ""
	for i := 0; i < len(res); i++ {
		if i != 0 {
			usersStr += " "
		}
		usersStr += res[i].Name
	}
	return mcp.NewToolResultText(usersStr), nil
}
