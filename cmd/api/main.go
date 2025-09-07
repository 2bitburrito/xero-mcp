package main

import (
	"context"
	"log"
	"net/http"

	mcpServer "github.com/2bitburrito/xero-mcp/internal/mcp"
	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	url := "localhost:8090"
	x := xeroapi.Xero{}
	err := x.Authorize()
	if err != nil {
		log.Fatalf("Couldn't Authorize Xero: %v", err)
	}
	ctx := context.Background()

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		server := mcpServer.NewServer(ctx, x)
		return server
	}, nil)

	log.Printf("MCP server listening on %s...", url)

	if err := http.ListenAndServe(url, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
