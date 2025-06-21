package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repository"
)

type Message struct {
	// メッセージUUID
	Id string `json:"id"`
	// 投稿者名
	UserName string `json:"userName"`
	// チャンネルUUID
	ChannelId string `json:"channelId"`
	// メッセージ本文
	Content string `json:"content"`
	// 投稿日時
	CreatedAt time.Time `json:"createdAt"`
	// 押されているスタンプの配列
	Stamps []Stamp `json:"stamps"`
}
type Stamp struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func GetMessagesTool() mcp.Tool {
	tool := mcp.NewTool("messages",
		mcp.WithDescription("Get some messages"),
		mcp.WithString("channel",
			mcp.Description("channel name to get messages from"),
		),
	)
	return tool
}

func GetMessagesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traq_client := clients.GetTraqClient()
	channel, err := request.RequireString("channel")

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	channel_to_id, err := repository.GetChannelId(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id, ok := channel_to_id[channel]
	if !ok {
		return mcp.NewToolResultText("Channel not found"), nil
	}

	getReq := traq_client.ChannelApi.GetMessages(ctx, id)

	res, _, err := getReq.Execute()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	//stampsを編集
	stampNameMap, err := repository.GetIdToStamp(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	channelNameMap, err := repository.GetChannelName(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	userNameMap, err := repository.GetUserName(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	allMessages := []Message{}
	stampCountMap := make(map[string]int)
	const getLimit = 20             //取得するメッセージの上限
	for _, message_v := range res { //投稿ごと
		for _, v := range message_v.Stamps { //スタンプごと
			stampName := stampNameMap[v.StampId]
			stampCountMap[stampName]++
		}
		var message Message
		message.Id = message_v.Id
		message.Content = message_v.Content
		message.UserName = userNameMap[message_v.UserId]
		message.CreatedAt = message_v.CreatedAt
		message.ChannelId = channelNameMap[message_v.ChannelId]

		for k, v := range stampCountMap {
			stamp := Stamp{
				Name:  k,
				Count: v,
			}
			message.Stamps = append(message.Stamps, stamp)
		}
		allMessages = append(allMessages, message)
		if len(allMessages) >= getLimit {
			break
		}
	}

	jsonBytes, err := json.Marshal(allMessages)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
