package handlers

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func SerchTool() mcp.Tool {
	tool := mcp.NewTool("search",
		mcp.WithDescription("Search"),
		mcp.WithString("word",
			mcp.Required(),
			mcp.Description("search word"),
		),
		mcp.WithString("to",
			mcp.Required(),
			mcp.Description("username whose mentioned"),
		),
		mcp.WithString("from",
			mcp.Required(),
			mcp.Description("username that sent message"),
		),
		mcp.WithString("bot",
			mcp.Required(),
			mcp.Description("whether it is bot"),
		),
	)
	return tool
}
