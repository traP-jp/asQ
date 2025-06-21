package handlers

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func GetAllUsrs() mcp.Tool {
	tool := mcp.NewTool("getAllUsers",
		mcp.WithDescription("Get all traP users"),
	)
	return tool
}
