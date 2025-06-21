package handlers

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func GetUser() mcp.Tool {
	tool := mcp.NewTool("getUser",
		mcp.WithDescription("Get user"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("user name"),
		),
	)
	return tool
}
