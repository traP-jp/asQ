package handlers

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repositry"
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
	name_to_id, err := repositry.GetUserToId(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id, ok := name_to_id[name]
	if !ok {
		return mcp.NewToolResultText("Username not found"), nil
	}
	getReq := traq_client.UserApi.GetUser(ctx, id)

	//searchReq := traq_client.UserApi.GetUser(ctx)

	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
