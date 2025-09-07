package mcpServer

import (
	"context"

	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type XeroToolHandler struct {
	context    context.Context
	xeroClient *xeroapi.Xero
}

func NewServer(ctx context.Context, x xeroapi.Xero) *mcp.Server {
	xh := &XeroToolHandler{
		context:    ctx,
		xeroClient: &x,
	}
	opts := mcp.ServerOptions{}
	server := mcp.NewServer(&mcp.Implementation{Name: "xero-mcp", Version: "v0.1.0"}, &opts)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-items",
		Description: "This will list all available items",
	},
		xh.listItems)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-invoices",
		Description: "Returns a list of all invoices",
	},
		xh.listAllInvoices)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-contacts",
		Description: "Lists all available contacts from xero",
	},
		xh.getContacts)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create-invoice",
		Description: "This will list all available pre-made items",
	},
		xh.createInvoice)

	return server
}
