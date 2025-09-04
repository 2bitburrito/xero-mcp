package xeroapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/2bitburrito/xero-mcp/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

const (
	XeroURL        = "https://api.xero.com/api.xro/2.0/"
	baseAuthURL    = `https://login.xero.com/identity/connect/authorize?response_type=code&client_id=%s&redirect_uri=%s:5678/callback&scope=openid offline_access openid profile email accounting.transactions&state=%s`
	xeroConnectURL = "https://api.xero.com/connections"
)

type Auth struct {
	ClientID        string
	ClientSecret    string
	URL             string
	stateCode       string
	redirectURL     string
	callback        authCallback
	baseCallbackURI string
	Tokens          Tokens
	jwt             accessToken
}

type authCallback struct {
	code         string
	state        string
	sessionState string
}
type tokenReq struct {
	GrantType   string `json:"grant_type"`
	RedirectURI string `json:"redirect_uri"`
	Code        string `json:"code"`
}

func (x *Xero) Authorize() error {
	x.Auth.stateCode = rand.Text()
	x.Auth.redirectURL = fmt.Sprintf(x.Auth.URL, x.Auth.ClientID, x.Auth.baseCallbackURI, x.Auth.stateCode)
	if err := utils.OpenURL(x.Auth.redirectURL); err != nil {
		return err
	}
	fmt.Println("redirect url: ", x.Auth.redirectURL)

	callbackChan := make(chan authCallback)
	errChan := make(chan error)

	// I know this is redundant - I wanted to practice writing
	// goroutines using channels
	go x.handleCallback(callbackChan, errChan)

	select {
	case cb := <-callbackChan:
		fmt.Printf("Received Callback: %+v\n", cb)
		x.Auth.callback = cb
		if err := x.getBearerToken(); err != nil {
			return err
		}
	case err := <-errChan:
		return fmt.Errorf("error while handling auth callback: %v", err)
	}
	id, err := x.getTennantID()
	return nil
}

func (x *Xero) handleCallback(callbackChan chan authCallback, errChan chan error) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		rCode := params.Get("code")
		rState := params.Get("state")
		rSessionState := params.Get("session_state")
		if rState != x.Auth.stateCode {
			errChan <- fmt.Errorf("received incorrect state from callback. Wanted: %s, received: %s",
				x.Auth.stateCode, rState)
			return
		}
		if len(rCode) == 0 || len(rSessionState) == 0 {
			errChan <- fmt.Errorf("didn't receive expected values while handling callback")
			return
		}
		callbackChan <- authCallback{
			code:         rCode,
			state:        rState,
			sessionState: rSessionState,
		}
	})
	addr := fmt.Sprintf(":%d", x.port)
	fmt.Printf("Listening for callback at %s/callback\n", addr)
	http.ListenAndServe(addr, nil)
}

func (x *Xero) getBearerToken() error {
	tokenURL := "https://identity.xero.com/connect/token"
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", x.Auth.callback.code)
	redirectURI := fmt.Sprintf("%s:%d/callback", x.Auth.baseCallbackURI, x.port)
	formData.Set("redirect_uri", redirectURI)
	encodedData := formData.Encode()

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(encodedData))
	if err != nil {
		return err
	}
	authHeaderRaw := fmt.Sprintf("%s:%s", x.Auth.ClientID, x.Auth.ClientSecret)
	authHeaderB64 := base64.StdEncoding.EncodeToString([]byte(authHeaderRaw))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+authHeaderB64)

	resp, err := x.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad response from token request: %s - %s", resp.Status, string(body))
	}

	if err := json.Unmarshal(body, &x.Auth.Tokens); err != nil {
		return err
	}
	return nil
}

func (x *Xero) getTennantID() (string, error) {
	if len(x.TennantID) != 0 {
		return x.TennantID, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	urlWithParam := xeroConnectURL + "?authEventId=" + x.Auth.Tokens.AccessToken
	req, err := http.NewRequest("GET", urlWithParam, nil)
	if err != nil {
		return "", err
	}
	resp, err := x.client.Do(req)
	if err != nil {
		return "", err
	}
}
