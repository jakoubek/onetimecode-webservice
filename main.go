package main

import (
	"encoding/json"
	"fmt"
	"github.com/jakoubek/onetimecode-webservice/handler"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jakoubek/onetimecode-webservice/requestlogger"

	"github.com/gorilla/mux"
)

var starttime time.Time
var requests int64
var requestsOld int64

func main() {
	starttime = time.Now()
	loadRequestsFromFile()
	initLogWriter()

	r := mux.NewRouter()
	r.HandleFunc("/", rootInfo).Methods("GET")
	r.HandleFunc("/number", handler.Number).Methods("GET")
	r.HandleFunc("/alphanumeric", handler.Alphanumeric).Methods("GET")
	r.HandleFunc("/ksuid", handler.Ksuid).Methods("GET")
	r.HandleFunc("/uuid", handler.Uuid).Methods("GET")
	//r.HandleFunc("/random", processRandom).Methods("GET")
	r.HandleFunc("/status", processStatus).Methods("GET")
	log.Print("Starting server on " + getServerPort())
	http.ListenAndServe(getServerPort(), r)
}

func initLogWriter() {
	go func() {
		for true {
			if requests > requestsOld {
				requestlogger.SaveCounterfile(getCounterfile(), requests)
				requestsOld = requests
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func loadRequestsFromFile() {
	requests = requestlogger.ReadCounterfile(getCounterfile())
	requestsOld = requests
}

func logRequest() {
	requests++
}

func rootInfo(w http.ResponseWriter, r *http.Request) {

	type result struct {
		Result string `json:"result"`
		Info   string `json:"info"`
	}

	response := result{
		Result: "OK",
		Info:   "Go to https://www.onetimecode.net for information on how to access the API. See /status for API health.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func processStatus(w http.ResponseWriter, r *http.Request) {

	type answer struct {
		Result        string    `json:"result"`
		Info          string    `json:"info"`
		ServerStarted time.Time `json:"server_started"`
		Timestamp     int64     `json:"timestamp"`
		Requests      int64     `json:"requests"`
	}

	result := answer{
		Result:        "OK",
		Info:          "API fully operational",
		ServerStarted: starttime,
		Timestamp:     time.Now().Unix(),
		Requests:      requests,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}

func getCounterfile() string {
	if filename, ok := os.LookupEnv("COUNTERFILE"); ok {
		return filename
	}
	return "counter.json"
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
