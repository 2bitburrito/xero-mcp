// Package server is the entrypoint for the main server
package server

import (
	"encoding/json"
	"log"
	"net/http"

	mcpServer "github.com/2bitburrito/xero-mcp/internal/mcp"
	"github.com/2bitburrito/xero-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	xh := &mcpServer.XeroToolHandler{
		Context:    s.Ctx,
		XeroClient: s.Xero,
	}
	// Register routes
	mcpHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		server := mcpServer.NewServer(xh)
		return server
	}, nil)

	mux.HandleFunc("/callback", http.HandlerFunc(s.HandleAuthCallback))
	mux.HandleFunc("/health", http.HandlerFunc(s.CheckHealth))
	mux.HandleFunc("/success", http.HandlerFunc(s.serveSuccess))
	mux.HandleFunc("/hello", http.HandlerFunc(s.HelloWorldHandler))
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("MCP server is active. Use POST to start a session."))
			return
		}
		mcpHandler.ServeHTTP(w, r)
	})
	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) serveSuccess(w http.ResponseWriter, r *http.Request) {
	html := "<h1>SUCCESS</h1>"
	if _, err := w.Write([]byte(html)); err != nil {
		utils.ReturnJsonError(w, err, http.StatusInternalServerError, "Server couldn't serve success url")
	}
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Accept", "text/event-stream")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
