package mcpServer

import (
	"context"
	"fmt"

	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type (
	GetContactsParams  struct{}
	ListItemsParams    struct{}
	ListInvoicesParams struct{}
)

func (h *XeroToolHandler) listItems(ctx context.Context, req *mcp.CallToolRequest, _ *ListItemsParams) (*mcp.CallToolResult, []xeroapi.Item, error) {
	items, err := h.xeroClient.GetItems()
	if err != nil {
		return nil, []xeroapi.Item{}, fmt.Errorf("failed to get items: %w", err)
	}
	return nil, items, nil
}

func (h *XeroToolHandler) getContacts(ctx context.Context, req *mcp.CallToolRequest, _ *GetContactsParams) (*mcp.CallToolResult, []xeroapi.Contact, error) {
	contacts, err := h.xeroClient.GetContacts()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	return nil, contacts, nil
}

func (h *XeroToolHandler) listAllInvoices(ctx context.Context, req *mcp.CallToolRequest, _ *ListInvoicesParams) (*mcp.CallToolResult, []xeroapi.Invoice, error) {
	invoices, err := h.xeroClient.GetInvoices()
	if err != nil {
		return nil, []xeroapi.Invoice{}, fmt.Errorf("failed to get invoices: %w", err)
	}
	return nil, invoices, nil
}

func (h *XeroToolHandler) createInvoice(ctx context.Context, req *mcp.CallToolRequest, _ *ListInvoicesParams) (*mcp.CallToolResult, []xeroapi.Invoice, error) {
	invoices, err := h.xeroClient.GetInvoices()
	if err != nil {
		return nil, []xeroapi.Invoice{}, fmt.Errorf("failed to create invoice: %w", err)
	}
	return nil, invoices, nil
}
