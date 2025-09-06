package mcpServer

import (
	"log"

	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
)

func listItems(x xeroapi.Xero) []xeroapi.Item {
	items, err := x.GetItems()
	log.Fatalf("%w", err)
	return items
}

func listAllInvoices(x xeroapi.Xero) []xeroapi.Invoice {
	invoices, err := x.GetInvoices()
	log.Fatalf("%w", err)
	return invoices
}
