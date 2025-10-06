package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/2bitburrito/xero-mcp/internal/server"
	"github.com/2bitburrito/xero-mcp/internal/setup"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dependencies, err := setup.Dependencies()
	if err != nil {
		log.Fatal("error setting up dependencies: ", err)
	}
	server := server.NewServer(dependencies)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
