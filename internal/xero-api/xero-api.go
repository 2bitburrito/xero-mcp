package xeroapi

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

const (
	XeroURL = "https://api.xero.com/api.xro/2.0"
	AuthURL = `https://login.xero.com/identity/connect/authorize?response_type=code&client_id=%s&redirect_uri=http://localhost:5678/callback&scope=openid offline_access openid profile email accounting.transactions&state=&s`
)

type Xero struct {
	Url      string
	client   *http.Client
	ClientID string
	port     int
	Auth     Auth
}
type Auth struct {
	URL       string
	stateCode string
	callback  authCallback
}

type authCallback struct {
	code         string
	state        string
	sessionState string
}

func (x *Xero) Authorize() error {
	x.Auth.stateCode = rand.Text()

	authURL := fmt.Sprintf(x.Auth.URL, x.ClientID, x.Auth.stateCode)
	fmt.Println("Visit this URL to authorize: ", authURL)
	callbackChan := make(chan authCallback)

	go x.handleCallback(callbackChan)
	<-callbackChan
	return nil
}

func (x *Xero) handleCallback(callbackChan chan authCallback) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recieved callback")
		fmt.Println(r.Body)
		r.Body.Close()
		rCode := r.Header.Get("code")
		rState := r.Header.Get("state")
		rSessionState := r.Header.Get("session_state")
		if rState != x.Auth.stateCode {
			fmt.Errorf("received incorrect state from callback. Wanted: %s, received: %s",
				x.Auth.stateCode, rState)
			return
		}
		if len(rCode) == 0 || len(rSessionState) == 0 {
			fmt.Errorf("didn't receive expected values while handling callback")
			return
		}
		callbackChan <- authCallback{
			code:         rCode,
			state:        rState,
			sessionState: rSessionState,
		}
	})
	addr := fmt.Sprintf(":%d", x.port)
	fmt.Printf("Listening for callback on %s/callback\n", addr)
	http.ListenAndServe(addr, nil)
}

func (x *Xero) GetItems() {
}
