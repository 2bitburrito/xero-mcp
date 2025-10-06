package xeroapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (x *Xero) makeAPICall(reqType, path string, body any) (*http.Response, error) {
	url := XeroURL + path
	fmt.Println("making api call to: ", url)
	var reqbody io.Reader
	if body != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
		reqbody = buf
	}
	req, err := http.NewRequest(reqType, url, reqbody)
	if err != nil {
		err = fmt.Errorf("new request error for call to %s: %w", url, err)
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+x.Auth.Tokens.AccessToken)
	req.Header.Set("Xero-tenant-id", x.Auth.Tennants.TenantID)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := x.Client.Do(req)
	if err != nil {
		err = fmt.Errorf("new request error for call to %s: %w", url, err)
		log.Println(err)
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errResp map[string]any
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		defer resp.Body.Close()
		if err != nil {
			err = fmt.Errorf("error in request call to: %s\nCouldn't decode to json, : %w", url, err)
			log.Println(err)
			return nil, err
		}

		return nil, fmt.Errorf("error in request call to: %s\nResponded with: %+v", url, errResp)
	}

	return resp, nil
}
