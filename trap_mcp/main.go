package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/handlers"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"MCP server to acquire information about traP",
		"0.1.0",
	)

	mcpServer.AddTool(handlers.SearchTool(), handlers.TraqSearchHandler)
	mcpServer.AddTool(handlers.GetAllUsrsTool(), handlers.GetAllUsersHandler)
	mcpServer.AddTool(handlers.GetUserTool(), handlers.GetUserHandler)
	mcpServer.AddTool(handlers.GetMessagesTool(), handlers.GetMessagesHandler)

	if err := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			return clients.GetTraqContext()
		}),
	).Start(":8000"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
