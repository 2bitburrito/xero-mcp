package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonErrReturn struct {
	Error string `json:"error"`
}

func ReturnJsonError(w http.ResponseWriter, e error, statusCode int, msg ...string) {
	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	log.Println(e)
	rtnMap := JsonErrReturn{
		Error: e.Error(),
	}
	RespondWithJSON(w, statusCode, rtnMap)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
