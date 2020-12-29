package handler

import (
	"encoding/json"
	"github.com/jakoubek/onetimecode"
	"net/http"
	"strconv"
)

func Alphanumeric(w http.ResponseWriter, r *http.Request) {

	type answer struct {
		Result string `json:"result"`
	}

	q := r.URL.Query()

	format := q.Get("format")

	length := -1
	lengthStr := q.Get("length")
	if lengthStr != "" {
		length, _ = strconv.Atoi(lengthStr)
	}

	caseStr := q.Get("case")

	if format == "" {
		format = "json"
	}

	var code *onetimecode.Onetimecode
	if length > -1 {
		switch caseStr {
		case "upper":
			code = onetimecode.NewAlphanumericalCode(
				onetimecode.WithLength(length),
				onetimecode.WithUpperCase(),
			)
		case "lower":
			code = onetimecode.NewAlphanumericalCode(
				onetimecode.WithLength(length),
				onetimecode.WithLowerCase(),
			)
		case "":
			code = onetimecode.NewAlphanumericalCode(
				onetimecode.WithLength(length),
			)
		}
	} else {
		switch caseStr {
		case "upper":
			code = onetimecode.NewAlphanumericalCode(
				onetimecode.WithUpperCase(),
			)
		case "lower":
			code = onetimecode.NewAlphanumericalCode(
				onetimecode.WithLowerCase(),
			)
		case "":
			code = onetimecode.NewAlphanumericalCode()
		}
	}

	var result = answer{
		Result: code.Code(),
	}

	if format == "txt" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(code.Code()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}

}
