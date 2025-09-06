package xeroapi

import (
	"encoding/json"
	"fmt"
)

type ContactsResponse struct {
	Contacts []Contact `json:"Contacts"`
}

func (x *Xero) GetContacts() ([]Contact, error) {
	var contacts ContactsResponse
	resp, err := x.makeApiCall("GET", "Contacts", nil)
	if err != nil {
		return []Contact{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&contacts)
	if err != nil {
		return []Contact{}, fmt.Errorf("failed to decode contacts response: %w", err)
	}

	return contacts.Contacts, nil
}

func (x *Xero) GetContact(contactID string) (*Contact, error) {
	var contacts ContactsResponse
	path := fmt.Sprintf("Contacts/%s", contactID)
	resp, err := x.makeApiCall("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&contacts)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contact response: %w", err)
	}

	if len(contacts.Contacts) == 0 {
		return nil, fmt.Errorf("contact not found")
	}

	return &contacts.Contacts[0], nil
}
