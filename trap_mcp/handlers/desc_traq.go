package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func TraQDescTool() mcp.Tool {
	tool := mcp.NewTool("desc_traq",
		mcp.WithDescription("Introduction of traQ"),
	)
	return tool
}

func TraQDescHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content := "traQ is a communication tool service available only to members of a club called traP.Users can Efficiently share information in real time.Organize information through channels, real-time chat, file sharing, and integration with external tools.Every user can post messages on each channel as they wish, and these messages can be viewed by everyone.traQ also has a stamp feature, which allows you to stamp each message.Stamps are a very important feature for streamlining communication and conveying emotions and reactions without sending messages.In addition, club members can send direct messages to each other and reply to messages.BOTs also exist in traQ, and there is a wide range of BOTs that manage tasks and respond to specific words.Channels have a parent-child relationship and are categorized by channel topic.Members of the club have the authority to create channels, depending on the topic they want to talk about.gps/times/ has a home channel for each member of the club.traQ is very important in traP's circle activities and provides access to various services such as KnoQ, anke-to, and HackMD.Most of the circle communications for events are announced through traP.There are eight parent channels at the highest level of the hierarchy: events, general, gps, projects, random, services, team, and univ. All channels are derived from these channels.There are various topics, not only those directly related to circle activities, but also channels that introduce recommended games, channels that show live classes, and so on." //ここにMDの内容を入れる
	return mcp.NewToolResultText(content), nil
}
