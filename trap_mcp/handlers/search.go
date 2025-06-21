package handlers

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repository/id_to_channel"
	"github.com/traP-jp/h25s_05/trap_mcp/repository/id_to_user"
)

type SearchFoundMessage struct {
	User    string `json:"user"`
	Content string `json:"content"`
	Channel string `json:"channel"`
}

func SearchTool() mcp.Tool {
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

	foundMessages := make([]SearchFoundMessage, 0)
	idToUser, err := id_to_user.GetIdToUserId(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	idToChannel, err := id_to_channel.GetIdToChannel(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	for _, hit := range res.Hits {

		foundMessages = append(foundMessages, SearchFoundMessage{
			User:    idToUser[hit.GetUserId()],
			Content: hit.GetContent(),
			Channel: idToChannel[hit.GetChannelId()],
		})
	}
	jsonData, err := json.Marshal(foundMessages)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}
