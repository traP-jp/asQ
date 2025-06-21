package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func SerchTool() mcp.Tool {
	tool := mcp.NewTool("search",
		mcp.WithDescription("Search"),
		mcp.WithString("word",
			mcp.Description("search word"),
		),
		mcp.WithString("to",
			mcp.Description("username whose mentioned"),
		),
		mcp.WithString("from",
			mcp.Description("username that sent message"),
		),
		mcp.WithString("bot",
			mcp.Description("whether it is bot"),
		),
	)
	return tool
}

func TraqSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	users, _, err := traq_client.MessageApi.SearchMessages(ctx)
	word := request.
	name, err := request("word")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultResource(fmt.Sprintf("Hello, %s!", name)), nil
}
