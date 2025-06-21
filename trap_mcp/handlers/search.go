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
	searchReq := traq_client.MessageApi.SearchMessages(ctx)
	word, err := request.RequireString("word")
	if err == nil {
		fmt.Println(word)
		searchReq = searchReq.Word(word)
	}
	from, err := request.RequireStringSlice("from")
	if err == nil {
		searchReq = searchReq.From(from)
	}
	to, err := request.RequireStringSlice("to")
	if err == nil {
		searchReq = searchReq.To(to)
	}
	bot, err := request.RequireBool("bot")
	if err == nil {
		searchReq = searchReq.Bot(bot)
	}

	res, _, err := searchReq.Execute()

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(res.Hits[0].Content), nil
}
