package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
)

var (
	version      string
	buildVersion string
	buildTime    string
	fullCommit   string
	isDebugMode  string = "false"
)

func main() {

	s := NewServer("Onetimecode-API 1.0", getCounterfile())

	s.logger.Printf("Server Version %s is starting on %s...", buildVersion, getServerPort())
	s.logger.Printf("Counter file: %s...", getCounterfile())

	s.setupRoutes()

	http.ListenAndServe(getServerPort(), s.router)

}

func (s *server) setupRoutes() {

	svcRoutes := s.router.Methods("GET").Subrouter()
	svcRoutes.HandleFunc("/number", s.logRequest(s.handleNumber("json")))
	svcRoutes.HandleFunc("/number.txt", s.logRequest(s.handleNumber("txt")))
	svcRoutes.HandleFunc("/alphanumeric", s.logRequest(s.handleAlphanumeric("json")))
	svcRoutes.HandleFunc("/alphanumeric.txt", s.logRequest(s.handleAlphanumeric("txt")))
	svcRoutes.HandleFunc("/ksuid", s.logRequest(s.handleKsuid("json")))
	svcRoutes.HandleFunc("/ksuid.txt", s.logRequest(s.handleKsuid("txt")))
	svcRoutes.HandleFunc("/uuid", s.logRequest(s.handleUuid("json")))
	svcRoutes.HandleFunc("/uuid.txt", s.logRequest(s.handleUuid("txt")))
	svcRoutes.HandleFunc("/dice", s.logRequest(s.handleDice("json")))
	svcRoutes.HandleFunc("/dice.txt", s.logRequest(s.handleDice("txt")))
	svcRoutes.HandleFunc("/coin", s.logRequest(s.handleCoin("json")))
	svcRoutes.HandleFunc("/coin.txt", s.logRequest(s.handleCoin("txt")))
	svcRoutes.Use(LogRequestMiddleware)

	baseRoutes := s.router.Methods("GET").Subrouter()
	baseRoutes.HandleFunc("/", s.logRequest(s.handleIndex()))
	baseRoutes.HandleFunc("/status", s.logRequest(s.handleStatus("")))
	baseRoutes.HandleFunc("/status.txt", s.logRequest(s.handleStatus("txt")))
	baseRoutes.HandleFunc("/healthz", s.handleHealthz())
	baseRoutes.NotFoundHandler = s.handleNotFound()

}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.onetimecode.net/?ref=apiroot", http.StatusFound)
	}
}

func (s *server) handleStatus(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Result:        "OK",
		//Info:          "API fully operational",

		if format == "txt" {
			sort.SliceStable(s.logInfo.Routes, func(i, j int) bool {
				return s.logInfo.Routes[i].RouteName < s.logInfo.Routes[j].RouteName
			})
			var response string
			response = "STATUS\n"
			response += fmt.Sprintf("- Server started: %s\n", s.logInfo.ServerStartedAt.Format("2006-01-02 15:04:06"))
			response += fmt.Sprintf("- Requests      : %3d\n", s.logInfo.Requests)
			response += fmt.Sprintf("- Last request  : %s\n", s.logInfo.LastRequestAt.Format("2006-01-02 15:04:06"))
			response += "\nROUTES\n"
			for _, r := range s.logInfo.Routes {
				response += fmt.Sprintf(
					"- %-14s: %3d\n",
					r.RouteName,
					r.Requests,
				)
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(s.logInfo)
		}
	}
}

func (s *server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Result string `json:"result"`
			Info   string `json:"info"`
		}{
			Result: "OK",
			Info:   "API fully operational",
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
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
