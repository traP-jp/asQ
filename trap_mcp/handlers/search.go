package handlers

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func SerchTool() mcp.Tool {
	tool := mcp.NewTool("search",
		mcp.WithDescription("Search word"),
		mcp.WithString("word",
			mcp.Required(),
			mcp.Description("search word"),
		),
	)
	return tool
}
