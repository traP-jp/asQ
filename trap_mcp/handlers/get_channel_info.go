package handlers

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repository/id_to_channel"
)

func GetChannelInfoTool() mcp.Tool {
	tool := mcp.NewTool("getChannelInfo",
		mcp.WithDescription("Get channel info"),
		mcp.WithString("channelName",
			mcp.Description("Channel name to get details"),
		),
	)
	return tool
}

type Channel struct {
	Name     string   `json:"name"`
	Topic    string   `json:"topic"`
	Children []string `json:"children"`
}

func GetChannelInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idToChannel, err := id_to_channel.GetIdToChannel(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	channelName, err := request.RequireString("channelName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	channelId := ""
	for k, v := range idToChannel {
		if channelName == v {
			channelId = k
			break
		}
	}
	traq_client := clients.GetTraqClient()
	getReq := traq_client.ChannelApi.GetChannel(ctx, channelId)

	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	childName := []string{}
	for _, v := range res.GetChildren() {
		childName = append(childName, idToChannel[v])
	}
	channel := Channel{
		Name:     res.Name,
		Topic:    res.Topic,
		Children: childName,
	}

	jsonBytes, err := json.Marshal(channel)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
