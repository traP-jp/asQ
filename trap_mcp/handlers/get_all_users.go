package handlers

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func GetAllUsrsTool() mcp.Tool {
	tool := mcp.NewTool("getAllUsers",
		mcp.WithDescription("Show all traP users except Webhook. You can choose whether to include BOT users and suspected users."),
		mcp.WithBoolean("includeSuspended",
			mcp.Description("Include suspended users"),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("includeBot",
			mcp.Description("Include bot"),
			mcp.DefaultBool(false),
		),
	)
	return tool
}

func GetAllUsersHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	getReq := traq_client.UserApi.GetUsers(ctx)
	useSuspended, err := request.RequireBool("includeSuspended")
	if err == nil {
		getReq = getReq.IncludeSuspended(useSuspended)
	}
	useBot, err := request.RequireBool("includeBot")
	if err == nil {
		getReq = getReq.IncludeSuspended(useSuspended)
	}
	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	usersStr := ""
	for i := 0; i < len(res); i++ {
		if (!useBot && res[i].Bot) || (!useSuspended && res[i].State != 1) || strings.HasPrefix(res[i].Name, "Webhook#") {
			continue
		}
		if i != 0 {
			usersStr += " "
		}
		usersStr += res[i].Name
	}
	return mcp.NewToolResultText(usersStr), nil
}
