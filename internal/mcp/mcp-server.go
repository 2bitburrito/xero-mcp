package mcpServer

import (
	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewServer(x xeroapi.Xero) *mcp.Server {
	opts := mcp.ServerOptions{}
	server := mcp.NewServer(&mcp.Implementation{Name: "xero-mcp", Version: "v0.1.0"}, &opts)

	mcp.AddTool(server, &mcp.Tool{Name: "list-items", Description: "This will list all available pre-made items"}, listItems(x))
	mcp.AddTool(server, &mcp.Tool{Name: "list-invoices", Description: "This will list all available pre-made items"}, listAllInvoices(x))
	mcp.AddTool(server, &mcp.Tool{Name: "create-invoice", Description: "This will list all available pre-made items"}, getItems)
	// Run the server over stdin/stdout, until the client disconnects
	return server
}
