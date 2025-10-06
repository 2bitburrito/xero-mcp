package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/2bitburrito/xero-mcp/internal/setup"
	xeroapi "github.com/2bitburrito/xero-mcp/internal/xero-api"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	Xero *xeroapi.Xero
	Ctx  context.Context
	port int
}

func NewServer(params *setup.Injection) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		Ctx:  params.Ctx,
		Xero: params.Xero,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Println("Server Started Successfully On Port: ", NewServer.port)
	return server
}
