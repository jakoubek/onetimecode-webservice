package main

import (
	"encoding/json"
	"github.com/jakoubek/onetimecode"
	"net/http"
	"strconv"
)

func (s *server) handleNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	type answer struct {
		Result int64 `json:"result"`
	}

	q := r.URL.Query()

	format := q.Get("format")

	length := -1
	lengthStr := q.Get("length")
	if lengthStr != "" {
		length, _ = strconv.Atoi(lengthStr)
	}

	var min, max int
	minStr := q.Get("min")
	maxStr := q.Get("max")
	if minStr != "" {
		min, _ = strconv.Atoi(minStr)
	}
	if maxStr != "" {
		max, _ = strconv.Atoi(maxStr)
	}

	if format == "" {
		format = "json"
	}

	var code *onetimecode.Onetimecode
	if length > -1 {
		code = onetimecode.NewNumericalCode(
			onetimecode.WithLength(length),
		)
	} else {
		if min > 0 && max == 0 {
			code = onetimecode.NewNumericalCode(
				onetimecode.WithMin(min),
			)
		} else if min == 0 && max > 0 {
			code = onetimecode.NewNumericalCode(
				onetimecode.WithMax(max),
			)
		} else if min > 0 && max > 0 {
			code = onetimecode.NewNumericalCode(
				onetimecode.WithMinMax(min, max),
			)
		} else {
			code = onetimecode.NewNumericalCode()
		}
	}

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