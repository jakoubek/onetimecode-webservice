package main

import (
	"encoding/json"
	"github.com/jakoubek/onetimecode-webservice/algorithm"
	"net/http"
)

func (s *server) handleKsuid(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type answer struct {
			Result string `json:"result"`
		}

		code := algorithm.NewKsuid()

		var result = answer{
			Result: code,
		}

		if format == "txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(code))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
		}

	}
}
