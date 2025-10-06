package xeroapi

import "net/http"

type Xero struct {
	Client *http.Client
	Auth   Auth
	Url    string
	Port   int
}

type Item struct {
	ItemID              string `json:"ItemID"`
	Code                string `json:"Code"`
	Description         string `json:"Description"`
	PurchaseDescription string `json:"PurchaseDescription"`
	UpdatedDateUTC      string `json:"UpdatedDateUTC"`
	PurchaseDetails     struct {
		UnitPrice   float64 `json:"UnitPrice"`
		AccountCode string  `json:"AccountCode"`
		TaxType     string  `json:"TaxType"`
	} `json:"PurchaseDetails"`
	SalesDetails struct {
		UnitPrice   float64 `json:"UnitPrice"`
		AccountCode string  `json:"AccountCode"`
		TaxType     string  `json:"TaxType"`
	} `json:"SalesDetails"`
	Name                      string  `json:"Name"`
	IsTrackedAsInventory      bool    `json:"IsTrackedAsInventory"`
	IsSold                    bool    `json:"IsSold"`
	IsPurchased               bool    `json:"IsPurchased"`
	InventoryAssetAccountCode string  `json:"InventoryAssetAccountCode,omitempty"`
	TotalCostPool             float64 `json:"TotalCostPool,omitempty"`
	QuantityOnHand            float64 `json:"QuantityOnHand,omitempty"`
}

type Invoice struct {
	Type    string `json:"Type"`
	Contact struct {
		ContactID     string `json:"ContactID"`
		ContactStatus string `json:"ContactStatus"`
		Name          string `json:"Name"`
		Addresses     []struct {
			AddressType  string `json:"AddressType"`
			AddressLine1 string `json:"AddressLine1,omitempty"`
			AddressLine2 string `json:"AddressLine2,omitempty"`
			City         string `json:"City,omitempty"`
			PostalCode   string `json:"PostalCode,omitempty"`
		} `json:"Addresses"`
		Phones []struct {
			PhoneType string `json:"PhoneType"`
		} `json:"Phones"`
		UpdatedDateUTC string `json:"UpdatedDateUTC"`
		IsSupplier     string `json:"IsSupplier"`
		IsCustomer     string `json:"IsCustomer"`
	} `json:"Contact"`
	Date            string             `json:"Date"`
	DateString      string             `json:"DateString"`
	DueDate         string             `json:"DueDate"`
	DueDateString   string             `json:"DueDateString"`
	Status          string             `json:"Status"`
	LineAmountTypes string             `json:"LineAmountTypes"`
	LineItems       []InvoiceLineItems `json:"LineItems"`
	SubTotal        string             `json:"SubTotal"`
	TotalTax        string             `json:"TotalTax"`
	Total           string             `json:"Total"`
	UpdatedDateUTC  string             `json:"UpdatedDateUTC"`
	CurrencyCode    string             `json:"CurrencyCode"`
	InvoiceID       string             `json:"InvoiceID"`
	InvoiceNumber   string             `json:"InvoiceNumber"`
	Payments        []struct {
		Date      string `json:"Date"`
		Amount    string `json:"Amount"`
		PaymentID string `json:"PaymentID"`
	} `json:"Payments"`
	AmountDue      string `json:"AmountDue"`
	AmountPaid     string `json:"AmountPaid"`
	AmountCredited string `json:"AmountCredited"`
}
type InvoiceLineItems struct {
	ItemCode    string `json:"ItemCode"`
	Description string `json:"Description"`
	Quantity    string `json:"Quantity"`
	UnitAmount  string `json:"UnitAmount"`
	TaxType     string `json:"TaxType"`
	TaxAmount   string `json:"TaxAmount"`
	LineAmount  string `json:"LineAmount"`
	AccountCode string `json:"AccountCode"`
	AccountID   string `json:"AccountId"`
	Item        struct {
		ItemID string `json:"ItemID"`
		Name   string `json:"Name"`
		Code   string `json:"Code"`
	} `json:"Item"`
	Tracking []struct {
		TrackingCategoryID string `json:"TrackingCategoryID"`
		Name               string `json:"Name"`
		Option             string `json:"Option"`
	} `json:"Tracking"`
	LineItemID string `json:"LineItemID"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"` // Used to call the API
	IdToken      string `json:"id_token"`     // The token containing user identity details (only returned if OpenID Connect scopes are requested).
	ExpiresIn    int    `json:"expires_in"`   // The amount of seconds until the access token expires.
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"` // The token used to refresh the access token once it has expired (only returned if the offline_access scope is requested).
}

type accessToken struct {
	Nbf                   int    `json:"nbf"`
	Exp                   int    `json:"exp"`
	Iss                   string `json:"iss"`
	Aud                   string `json:"aud"`
	ClientID              string `json:"client_id"`
	Sub                   string `json:"sub"`
	AuthTime              int    `json:"auth_time"`
	XeroUserid            string `json:"xero_userid"`
	GlobalSessionID       string `json:"global_session_id"`
	Jti                   string `json:"jti"`
	AuthenticationEventID string `json:"authentication_event_id"`
	Scope                 string `json:"scope"`
}

type Contact struct {
	ContactID     string `json:"ContactID"`
	ContactStatus string `json:"ContactStatus"`
	Name          string `json:"Name"`
	FirstName     string `json:"FirstName,omitempty"`
	LastName      string `json:"LastName,omitempty"`
	EmailAddress  string `json:"EmailAddress,omitempty"`
	Addresses     []struct {
		AddressType  string `json:"AddressType"`
		AddressLine1 string `json:"AddressLine1,omitempty"`
		AddressLine2 string `json:"AddressLine2,omitempty"`
		City         string `json:"City,omitempty"`
		PostalCode   string `json:"PostalCode,omitempty"`
		Country      string `json:"Country,omitempty"`
	} `json:"Addresses"`
	Phones []struct {
		PhoneType        string `json:"PhoneType"`
		PhoneNumber      string `json:"PhoneNumber,omitempty"`
		PhoneAreaCode    string `json:"PhoneAreaCode,omitempty"`
		PhoneCountryCode string `json:"PhoneCountryCode,omitempty"`
	} `json:"Phones"`
	UpdatedDateUTC string `json:"UpdatedDateUTC"`
	IsSupplier     bool   `json:"IsSupplier"`
	IsCustomer     bool   `json:"IsCustomer"`
}

type TennantResponse struct {
	ID             string `json:"id"`
	AuthEventID    string `json:"authEventId"`
	TenantID       string `json:"tenantId"`
	TenantType     string `json:"tenantType"`
	TenantName     string `json:"tenantName"`
	CreatedDateUtc string `json:"createdDateUtc"`
	UpdatedDateUtc string `json:"updatedDateUtc"`
}
