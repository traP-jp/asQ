package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"MCP server to acquire information about traP",
		"0.1.0",
	)

	if err := server.NewStreamableHTTPServer(mcpServer).Start(":8000"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
