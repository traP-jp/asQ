package main

import (
	"context"
	"fmt"
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

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Checking Authorization header")
			fmt.Println("Gotten Token:", r.Header.Get("Authorization"))
			fmt.Println("Expected Token:", "Bearer "+bearerToken)
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
