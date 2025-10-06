package xeroapi

import (
	"encoding/json"
	"fmt"
)

type InvoicesResponse struct {
	Invoices []Invoice `json:"Invoices"`
}

//	{
//	  "Type": "ACCREC",
//	  "Contact": {
//	    "ContactID": "eaa28f49-6028-4b6e-bb12-d8f6278073fc"
//	  },
//	  "Date": "\/Date(1518685950940+0000)\/",
//	  "DateString": "2009-05-27T00:00:00",
//	  "DueDate": "\/Date(1518685950940+0000)\/",
//	  "DueDateString": "2009-06-06T00:00:00",
//	  "LineAmountTypes": "Exclusive",
//	  "LineItems": [
//	    {
//	      "Description": "Consulting services as agreed (20% off standard rate)",
//	      "Quantity": "10",
//	      "UnitAmount": "100.00",
//	      "AccountCode": "200",
//	      "DiscountRate": "20"
//	    }
//	  ]
//	}
type CreateInvoiceRequest struct {
	Type    string `json:"Type"`
	Contact struct {
		ContactID string `json:"ContactID"`
	} `json:"Contact"`
	DueDate         string         `json:"DueDate"`
	LineAmountTypes string         `json:"LineAmountTypes"`
	LineItems       []InvoiceItems `json:"LineItems"`
}
type InvoiceItems struct {
	Description  string `json:"Description"`
	LineItemID   string `json:"LineItemID,omitempty"`
	Quantity     int    `json:"Quantity"`
	UnitAmount   int    `json:"UnitAmount"`
	DiscountRate string `json:"DiscountRate"`
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

type CreateInvoiceParams struct {
	ContactID string
	Items     []struct {
		Quantity    int
		LineItemID  string
		ItemCode    string
		Description string
	}
}

func (x *Xero) CreateInvoice(params CreateInvoiceParams) (*Invoice, error) {
	inv := CreateInvoiceRequest{
		Type:    "ACCREC",
		DueDate: Date.AddDate(0, 0, 30).Format("2006-01-02"),
		Contact: struct {
			ContactID string `json:"ContactID"`
		}{
			ContactID: params.ContactID,
		},
		LineAmountTypes: "Exclusive",
	}
	for _, item := range params.Items {
		newItem := InvoiceItems{
			Quantity:    item.Quantity,
			LineItemID:  item.LineItemID,
			Description: item.Description,
		}
		inv.LineItems = append(inv.LineItems, newItem)
	}
	resp, err := x.makeAPICall("POST", "Invoices", inv)
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
	// TODO: Fill in the fields to update

	requestBody := CreateInvoiceRequest{
		// Invoices: []Invoice{invoice},
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
