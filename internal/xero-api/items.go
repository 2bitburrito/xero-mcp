package xeroapi

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("RESP in getItems:", resp)
	err = json.NewDecoder(resp.Body).Decode(&items)
	fmt.Printf("ITEMS: %+v\n", items.Items)
	return items.Items, err
}
