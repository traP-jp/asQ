package main

import (
	"context"
	"net/http"
	"os"

	"github.com/traP-jp/h25s_05/trap_mcp/clients"
	"github.com/traP-jp/h25s_05/trap_mcp/handlers"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"MCP server to acquire information about traP",
		"0.1.0",
	)
	bearerToken := os.Getenv("MCP_SERVER_TOKEN")

	mcpServer.AddTool(handlers.SearchTool(), handlers.TraqSearchHandler)
	mcpServer.AddTool(handlers.GetAllUsrsTool(), handlers.GetAllUsersHandler)
	mcpServer.AddTool(handlers.GetUserTool(), handlers.GetUserHandler)
	mcpServer.AddTool(handlers.GetKnoqTool(), handlers.GetKnoqHandler)
	mcpServer.AddTool(handlers.GetMessageTool(), handlers.GetMessageHandler)
	mcpServer.AddTool(handlers.SearchMdTool(), handlers.SearchMdHandler)
	mcpServer.AddTool(handlers.GetAllChannelsTool(), handlers.GetAllChannelsHandler)
	mcpServer.AddTool(handlers.TraQDescTool(), handlers.TraQDescHandler)

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+bearerToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	streamableServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			return clients.GetTraqContext()
		}),
	)
	http.ListenAndServe(":8081", authMiddleware(streamableServer))
}
