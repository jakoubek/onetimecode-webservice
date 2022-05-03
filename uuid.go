package main

import (
	"encoding/json"
	"github.com/jakoubek/onetimecode"
	"net/http"
	"strconv"
)

func (s *server) handleUuid(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type answer struct {
			Result string `json:"result"`
		}

		q := r.URL.Query()

		withoutDashes, _ := strconv.ParseBool(q.Get("withoutdashes"))

		var code *onetimecode.Onetimecode
		if withoutDashes == true {
			code = onetimecode.NewUuidCode(
				onetimecode.WithoutDashes(),
			)
		} else {
			code = onetimecode.NewUuidCode()
		}

		var result = answer{
			Result: code.Code(),
		}

		if format == "txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(code.Code()))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
		}

	}
}
