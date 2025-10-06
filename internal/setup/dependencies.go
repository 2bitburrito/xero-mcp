// Package setup handles dependencies for the application are defined here
package setup

import (
	"context"
	"net/http"
	"os"
	"strconv"

	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	"github.com/joho/godotenv"
)

type Injection struct {
	Xero *xeroapi.Xero
	Ctx  context.Context
}

func Dependencies() (*Injection, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Port is nil in depencdencies")
	}
	return &Injection{
		Xero: &xeroapi.Xero{
			Port:   port,
			Client: &http.Client{},
			Url:    xeroapi.XeroURL,
			Auth: xeroapi.Auth{
				URL:             xeroapi.BaseAuthURL,
				ClientID:        os.Getenv("XERO_CLIENT_ID"),
				ClientSecret:    os.Getenv("XERO_CLIENT_SECRET"),
				BaseCallbackURI: os.Getenv("XERO_CLIENT_CALLBACK_URI"),
			},
		},
		Ctx: context.Background(),
	}, nil
}
