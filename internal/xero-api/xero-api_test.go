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
			URL: AuthURL,
		},
		client:   &http.Client{},
		ClientID: os.Getenv("XERO_CLIENT_ID"),
		port:     5678,
	}

	if err := xero.Authorize(); err != nil {
		t.Fail()
		fmt.Println(err)
	}
}
