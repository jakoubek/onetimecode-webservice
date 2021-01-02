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
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(code))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
		}

	}
}
