// Package mcpserver is a package that provides the Model Context Protocol (MCP) server
package mcpserver

import (
	"context"

	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type XeroToolHandler struct {
	Context    context.Context
	XeroClient *xeroapi.Xero
}

func NewServer(xh *XeroToolHandler) *mcp.Server {
	opts := mcp.ServerOptions{
		Instructions: "Ensure that you have run 'authenticate' before calling any other tools",
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "xero-mcp", Version: "v0.1.0"}, &opts)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "authenticate",
		Description: "This will return a url the user must click through to authenticate",
	},
		xh.getAuthURL)

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

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "create-invoice",
	// 	Description: "This will list all available pre-made items",
	// },
	// 	xh.createInvoice)

	return server
}
