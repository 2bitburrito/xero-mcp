package xeroapi

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/2bitburrito/xero-mcp/internal/utils"
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
			URL:             BaseAuthURL,
			ClientID:        os.Getenv("XERO_CLIENT_ID"),
			ClientSecret:    os.Getenv("XERO_CLIENT_SECRET"),
			BaseCallbackURI: os.Getenv("XERO_CLIENT_CALLBACK_URI"),
		},
		Port:   5678,
		Client: &http.Client{},
	}

	authURL := xero.GetAuthURL()
	if len(authURL) == 0 {
		fmt.Println("returned url of len 0")
		t.FailNow()
	}
	fmt.Println("redirect url: ", authURL)
	utils.OpenURL(authURL)

	wg := sync.WaitGroup{}
	wg.Add(1)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		params := r.URL.Query()
		if err := xero.HandleCallback(params); err != nil {
			fmt.Printf("error while handling auth callback: %v\n", err)
			t.FailNow()
		}
	})

	addr := fmt.Sprintf(":%d", xero.Port)
	fmt.Printf("Listening for callback at %s/callback\n", addr)

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println("error starting server: ", err)
		}
	}()
	wg.Wait()

	if err = xero.getTennantID(); err != nil {
		fmt.Println("error getting tennant ID: ", err)
		t.FailNow()
	}

	if err := xero.refreshJwt(); err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	items, err := xero.GetItems()
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
	fmt.Println(items)
}
