// Package xeroapi interacts with the xero api and handles auth
package xeroapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	XeroURL     = "https://api.xero.com/api.xro/2.0/"
	BaseAuthURL = `https://login.xero.com/identity/connect/authorize?response_type=code&client_id=%s&redirect_uri=%s:%v/callback&scope=openid offline_access openid profile email accounting.transactions accounting.settings&state=%s`
)

type Auth struct {
	ClientID        string
	ClientSecret    string
	URL             string
	stateCode       string
	redirectURL     string
	Callback        AuthCallback
	BaseCallbackURI string
	Tokens          TokenResp
	jwt             accessToken
	Tennants        TennantResponse
}

type AuthCallback struct {
	code         string
	state        string
	sessionState string
}

// GetAuthURL constructs the OAuth2 authorization URL and generates a random state code
func (x *Xero) GetAuthURL() string {
	x.Auth.stateCode = rand.Text()
	fmt.Println("setting stateCode: ", x.Auth.stateCode)
	x.Auth.redirectURL = fmt.Sprintf(x.Auth.URL, x.Auth.ClientID, x.Auth.BaseCallbackURI, x.Port, x.Auth.stateCode)
	return x.Auth.redirectURL
}

func (x *Xero) HandleCallback(params url.Values) error {
	rCode := params.Get("code")
	rState := params.Get("state")
	rSessionState := params.Get("session_state")
	fmt.Println("State code on struct: ", x.Auth.stateCode)
	if rState != x.Auth.stateCode {
		return fmt.Errorf("received incorrect state from callback. Received: %s, Wanted: %s",
			rState, x.Auth.stateCode)
	}
	if len(rCode) == 0 || len(rSessionState) == 0 {
		return fmt.Errorf("didn't receive expected values while handling callback")
	}
	x.Auth.Callback = AuthCallback{
		code:         rCode,
		state:        rState,
		sessionState: rSessionState,
	}
	if err := x.exchangeForBearerToken(); err != nil {
		return fmt.Errorf("error exchanging bearer tokens in callback")
	}
	if err := x.getTennantID(); err != nil {
		return fmt.Errorf("error getting tennant ID: %w", err)
	}
	return nil
}

func (x *Xero) exchangeForBearerToken() error {
	tokenURL := "https://identity.xero.com/connect/token"
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", x.Auth.Callback.code)
	redirectURI := fmt.Sprintf("%s:%d/callback", x.Auth.BaseCallbackURI, x.Port)
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

	resp, err := x.Client.Do(req)
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
	x.decodeJWT()
	return nil
}

func (x *Xero) decodeJWT() error {
	if len(x.Auth.Tokens.AccessToken) == 0 {
		return fmt.Errorf("auth Token not found in xero struct")
	}
	jwt := strings.Split(x.Auth.Tokens.AccessToken, ".")
	if len(jwt) < 2 {
		return fmt.Errorf("JWT is invalid or missing")
	}
	payload, _ := base64.RawURLEncoding.DecodeString(jwt[1])
	json.Unmarshal(payload, &x.Auth.jwt)
	return nil
}

// getTennantID retrieves the tenant ID associated with the authenticated user
// this is required for making API calls to Xero
func (x *Xero) getTennantID() error {
	xeroConnectURL := "https://api.xero.com/connections"

	if len(x.Auth.Tennants.TenantID) != 0 {
		return nil
	}
	req, err := http.NewRequest("GET", xeroConnectURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+x.Auth.Tokens.AccessToken)

	resp, err := x.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("responded with Error Code: %d", resp.StatusCode)
	}
	var tennant []TennantResponse
	err = json.NewDecoder(resp.Body).Decode(&tennant)
	if err != nil {
		return err
	}
	x.Auth.Tennants = tennant[0]
	return nil
}

func (x *Xero) refreshJwt() error {
	refreshURL := "https://identity.xero.com/connect/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", x.Auth.Tokens.RefreshToken)

	fmt.Println("encoded data in refreshJWT:", data.Encode())
	req, err := http.NewRequest("POST", refreshURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+base64.RawStdEncoding.EncodeToString([]byte(x.Auth.ClientID+":"+x.Auth.ClientSecret)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var respData TokenResp
	resp, err := x.Client.Do(req)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Println(resp)
		return fmt.Errorf("error in refresh_token: %+v", respData)
	}
	fmt.Printf("RESP in RefreshToken: %+v", respData)
	x.Auth.Tokens = respData
	err = x.decodeJWT()
	if err != nil {
		return err
	}

	return nil
}
