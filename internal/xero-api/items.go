package xeroapi

import (
	"encoding/json"
)

type ItemsResponse struct {
	Items []Item `json:"Items"`
}

func (x *Xero) GetItems() ([]Item, error) {
	var items ItemsResponse
	resp, err := x.makeApiCall("GET", "Items", nil)
	if err != nil {
		return []Item{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&items)
	return items.Items, err
}
