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

type EmailInvoiceRequest struct {
	To []string `json:"To,omitempty"`
}

func (x *Xero) GetInvoices() ([]Invoice, error) {
	var invoices InvoicesResponse
	resp, err := x.makeAPICall("GET", "Invoices", nil)
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
	resp, err := x.makeAPICall("GET", path, nil)
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

	resp, err := x.makeAPICall("POST", "Invoices", requestBody)
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

func (x *Xero) UpdateInvoice(invoiceID string, invoice Invoice) (*Invoice, error) {
	requestBody := CreateInvoiceRequest{
		Invoices: []Invoice{invoice},
	}

	path := fmt.Sprintf("Invoices/%s", invoiceID)
	resp, err := x.makeAPICall("POST", path, requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var invoicesResp InvoicesResponse
	err = json.NewDecoder(resp.Body).Decode(&invoicesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode update invoice response: %w", err)
	}

	if len(invoicesResp.Invoices) == 0 {
		return nil, fmt.Errorf("no invoice returned from update request")
	}

	return &invoicesResp.Invoices[0], nil
}

func (x *Xero) EmailInvoice(invoiceID string) error {
	path := fmt.Sprintf("Invoices/%s/Email", invoiceID)
	resp, err := x.makeAPICall("POST", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (x *Xero) EmailInvoiceTo(invoiceID string, emails []string) error {
	path := fmt.Sprintf("Invoices/%s/Email", invoiceID)
	requestBody := EmailInvoiceRequest{
		To: emails,
	}
	resp, err := x.makeAPICall("POST", path, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
