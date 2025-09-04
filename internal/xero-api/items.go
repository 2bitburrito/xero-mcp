package xeroapi

import (
	"encoding/json"
	"fmt"
)

func (x *Xero) GetItems() ([]Items, error) {
	var items []Items
	resp, err := x.makeApiCall("GET", "Items", nil)
	if err != nil {
		return []Items{}, err
	}
	fmt.Println(resp)
	err = json.NewDecoder(resp.Body).Decode(&items)
	return items, err
}
