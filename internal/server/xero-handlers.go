package server

import (
	"fmt"
	"net/http"

	"github.com/2bitburrito/xero-mcp/internal/utils"
)

func (s *Server) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	fmt.Println("HIT")
	if err := s.Xero.HandleCallback(params); err != nil {
		utils.ReturnJsonError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println("callback received and jwt set")
}
