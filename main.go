package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jakoubek/onetimecode"
)

type answer struct {
	Code   string `json:"code"`
	Mode   string `json:"mode"`
	Length int    `json:"length"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/onetime", processOnetimecode).Methods("GET")
	http.ListenAndServe(getServerPort(), r)
}

func processOnetimecode(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	mode := q.Get("mode")
	lengthStr := q.Get("length")

	if mode == "" {
		mode = "numbers"
	}

	var length int

	if lengthStr == "" {
		length = 6
	} else {
		length, _ = strconv.Atoi(lengthStr)
		if length <= 0 {
			length = 6
		}
	}

	var code string

	switch mode {
	case "numbers":
		code = onetimecode.NumberCode(length)
	case "alphanum":
		code = onetimecode.AlphaNumberCode(length)
	case "alphanumuc":
		code = onetimecode.AlphaNumberUcCode(length)
	default:
		mode = "numbers"
		code = onetimecode.NumberCode(length)
	}

	var result = answer{
		Code:   code,
		Mode:   mode,
		Length: length,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func getServerPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return ":3000"
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>We could not find the page you were looking for :(</h1><p>Please email us if you keep being sent to an "+
		"invalid page.</p>")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
