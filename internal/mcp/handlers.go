package mcpserver

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

func (h *XeroToolHandler) getAuthURL(ctx context.Context, req *mcp.CallToolRequest, _ *ListItemsParams) (*mcp.CallToolResult, map[string]string, error) {
	authURL := h.XeroClient.GetAuthURL()
	m := make(map[string]string)
	m["instructions"] = fmt.Sprintf("Please get the user to follow this link: \n\n`%s`", (authURL))
	m["auth_url"] = authURL
	fmt.Printf("sending authurl: %+v\n", m)
	return nil, m, nil
}

func (h *XeroToolHandler) listItems(ctx context.Context, req *mcp.CallToolRequest, _ *ListItemsParams) (*mcp.CallToolResult, map[string]xeroapi.Item, error) {
	m := make(map[string]xeroapi.Item)
	items, err := h.XeroClient.GetItems()
	if err != nil {
		return nil, m, fmt.Errorf("failed to get items: %w", err)
	}
	for _, item := range items {
		m[item.ItemID] = item
	}
	return nil, m, nil
}

func (h *XeroToolHandler) getContacts(ctx context.Context, req *mcp.CallToolRequest, _ *GetContactsParams) (*mcp.CallToolResult, map[string]xeroapi.Contact, error) {
	m := make(map[string]xeroapi.Contact)
	contacts, err := h.XeroClient.GetContacts()
	if err != nil {
		return nil, m, fmt.Errorf("failed to get contacts: %w", err)
	}
	for _, contact := range contacts {
		m[contact.ContactID] = contact
	}

	return nil, m, nil
}

func (h *XeroToolHandler) listAllInvoices(ctx context.Context, req *mcp.CallToolRequest, _ *ListInvoicesParams) (*mcp.CallToolResult, map[string]xeroapi.Invoice, error) {
	m := make(map[string]xeroapi.Invoice)
	invoices, err := h.XeroClient.GetInvoices()
	if err != nil {
		return nil, m, fmt.Errorf("failed to get invoices: %w", err)
	}
	for _, invoice := range invoices {
		m[invoice.InvoiceID] = invoice
	}
	return nil, m, nil
}

// func (h *XeroToolHandler) createInvoice(ctx context.Context, req *mcp.CallToolRequest, _ *ListInvoicesParams) (*mcp.CallToolResult, xeroapi.Invoice, error) {
// 	invoices, err := h.XeroClient.GetInvoices()
// 	if err != nil {
// 		return nil, xeroapi.Invoice{}, fmt.Errorf("failed to create invoice: %w", err)
// 	}
// 	return nil, invoices, nil
// }
