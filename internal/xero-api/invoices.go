package xeroapi

import (
	"encoding/json"
	"fmt"
)

type InvoicesResponse struct {
	Invoices []Invoice `json:"Invoices"`
}

type CreateInvoiceRequest struct {
	Invoices []Invoice `json:"Invoices"`
}

func (x *Xero) GetInvoices() ([]Invoice, error) {
	var invoices InvoicesResponse
	resp, err := x.makeApiCall("GET", "Invoices", nil)
	if err != nil {
		return []Invoice{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&invoices)
	if err != nil {
		return []Invoice{}, fmt.Errorf("failed to decode invoices response: %w", err)
	}

	return invoices.Invoices, nil
}

func (x *Xero) GetInvoice(invoiceID string) (*Invoice, error) {
	var invoices InvoicesResponse
	path := fmt.Sprintf("Invoices/%s", invoiceID)
	resp, err := x.makeApiCall("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&invoices)
	if err != nil {
		return nil, fmt.Errorf("failed to decode invoice response: %w", err)
	}

	if len(invoices.Invoices) == 0 {
		return nil, fmt.Errorf("invoice not found")
	}

	return &invoices.Invoices[0], nil
}

func (x *Xero) CreateInvoice(invoice Invoice) (*Invoice, error) {
	requestBody := CreateInvoiceRequest{
		Invoices: []Invoice{invoice},
	}

	resp, err := x.makeApiCall("POST", "Invoices", requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var invoicesResp InvoicesResponse
	err = json.NewDecoder(resp.Body).Decode(&invoicesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode create invoice response: %w", err)
	}

	if len(invoicesResp.Invoices) == 0 {
		return nil, fmt.Errorf("no invoice returned from create request")
	}

	return &invoicesResp.Invoices[0], nil
}
