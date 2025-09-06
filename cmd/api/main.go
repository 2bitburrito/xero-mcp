package main

import (
	"log"
	"net/http"

	mcpServer "github.com/2bitburrito/xero-mcp/internal/mcp"
	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	url := "http://localhost:9876"
	x := xeroapi.Xero{}

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		server := mcpServer.NewServer(x)
		return server
	}, nil)

	log.Printf("MCP server listening on %s", url)
	// Start the HTTP server with logging handler.
	if err := http.ListenAndServe(url, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
