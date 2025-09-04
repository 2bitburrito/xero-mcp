package xeroapi

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestXeroApi(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	xero := Xero{
		Url: XeroURL,
		Auth: Auth{
			URL:             baseAuthURL,
			ClientID:        os.Getenv("XERO_CLIENT_ID"),
			ClientSecret:    os.Getenv("XERO_CLIENT_SECRET"),
			baseCallbackURI: os.Getenv("XERO_CLIENT_CALLBACK_URI"),
		},
		client: &http.Client{},
		port:   5678,
	}

	if err := xero.Authorize(); err != nil {
		t.Fail()
		fmt.Println(err)
	}

	items, err := xero.GetItems()
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
	fmt.Println(items)
}
