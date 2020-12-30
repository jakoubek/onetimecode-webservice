package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main() {

	s := NewServer("Onetimecode-API 1.0", getCounterfile())

	s.logger.Printf("Server is starting on %s...", getServerPort())
	s.logger.Printf("Counter file: %s...", getCounterfile())

	s.setupRoutes()

	http.ListenAndServe(getServerPort(), s.router)
}

func (s *server) setupRoutes() {
	s.router.HandleFunc("/", s.logRequest(s.handleIndex()))
	s.router.HandleFunc("/status", s.logRequest(s.handleStatus()))
	s.router.HandleFunc("/number", s.logRequest(s.handleNumber()))
	s.router.HandleFunc("/alphanumeric", s.logRequest(s.handleAlphanumeric()))
	s.router.HandleFunc("/ksuid", s.logRequest(s.handleKsuid()))
	s.router.HandleFunc("/uuid", s.logRequest(s.handleUuid()))
	s.router.NotFoundHandler = s.handleNotFound()
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Result string `json:"result"`
			Info   string `json:"info"`
		}{
			Result: "OK",
			Info:   "Go to https://www.onetimecode.net for information on how to access the API. See /status for API health.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (s *server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

			//Result:        "OK",
			//Info:          "API fully operational",

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.logInfo)
	}
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
