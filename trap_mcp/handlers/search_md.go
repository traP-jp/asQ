package handlers

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	mdsearch "github.com/traP-jp/h25s_05/trap_mcp/repository/md_search"
)

func SearchMdTool() mcp.Tool {
	tool := mcp.NewTool("searchMarkdown",
		mcp.WithDescription("Search and find markdown documents by keyword"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("Keyword to search for in markdown documents"),
		),
	)
	return tool
}

func SearchMdHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keyword, err := request.RequireString("keyword")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	results, err := mdsearch.Search(ctx, keyword)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(resultsJSON)), nil
}
