package xeroapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (x *Xero) makeApiCall(reqType, path string, body any) (*http.Response, error) {
	url := XeroURL + path
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
		return nil, fmt.Errorf("new request error for call to %s: %w", url, err)
	}
	req.Header.Set("Authorization", "Bearer "+x.Auth.Tokens.AccessToken)
	req.Header.Set("Xero-tenant-id", x.Auth.Tennants.TenantID)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := x.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("new request error for call to %s: %w", url, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errResp map[string]any
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		defer resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error in request call to: %s\nCouldn't decode to json", url)
		}

		return nil, fmt.Errorf("error in request call to: %s\nResponded with: %+v", url, errResp)
	}

	return resp, nil
}
