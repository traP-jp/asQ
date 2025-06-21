package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/repository"
	"github.com/traPtitech/go-traq"
)

type Message struct {
	// メッセージUUID
	Id string `json:"id"`
	// 投稿者UUID
	UserId string `json:"userId"`
	// チャンネルUUID
	ChannelId string `json:"channelId"`
	// メッセージ本文
	Content string `json:"content"`
	// 投稿日時
	CreatedAt time.Time `json:"createdAt"`
	// 押されているスタンプの配列
	Stamps []traq.MessageStamp `json:"stamps"`
	// スレッドUUID
	ThreadId traq.NullableString `json:"threadId"`
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
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}
