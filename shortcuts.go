package main

import (
	"encoding/json"
	"github.com/jakoubek/onetimecode"
	"net/http"
	"strconv"
)

func (s *server) handleDice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type answer struct {
			Result int64 `json:"result"`
		}

		q := r.URL.Query()

		format := q.Get("format")

		if format == "" {
			format = "json"
		}

		code := onetimecode.NewNumericalCode(
			onetimecode.WithMinMax(1, 6),
		)

		var result = answer{
			Result: code.NumberCode(),
		}

		if format == "txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte((strconv.FormatInt(code.NumberCode(), 10))))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
		}

	}
}

func (s *server) handleCoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type answer struct {
			Result int    `json:"result"`
			Side   string `json:"side"`
		}

		q := r.URL.Query()

		format := q.Get("format")

		if format == "" {
			format = "json"
		}

		code := onetimecode.NewNumericalCode(
			onetimecode.WithMinMax(0, 1),
		)

		var coinResult string
		switch code.NumberCode() {
		case 0:
			coinResult = "heads"
		case 1:
			coinResult = "tails"
		}

		result := answer{
			Result: int(code.NumberCode()),
			Side:   coinResult,
		}

		if format == "txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte((strconv.FormatInt(code.NumberCode(), 10))))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
		}

	}
}
