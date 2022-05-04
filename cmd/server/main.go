package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type config struct {
	port            int
	env             string
	logfilePath     string
	counterfilePath string
}

type application struct {
	config config

	startupTime time.Time

	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Web server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.logfilePath, "logfile", "./logfile.log", "Path and name of logfile")
	flag.StringVar(&cfg.counterfilePath, "counterfile", "", "Path and name of JSON counterfile")

	flag.Parse()

	// Setup logger
	logfileF, err := os.OpenFile(cfg.logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logfileF.Close()
	logger := log.New(logfileF, "", log.Ldate|log.Ltime)
	if cfg.env == "development" {
		mw := io.MultiWriter(os.Stdout, logfileF)
		logger.SetOutput(mw)
	}

	app := &application{
		config:      cfg,
		startupTime: time.Now(),
		logger:      logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	logger.Printf("Server version %s (%s)", buildVersion, buildTime)
	err = srv.ListenAndServe()
	logger.Fatal(err)

}

func (s *server) setupRoutes() {

	//svcRoutes := s.router.Methods("GET").Subrouter()
	//svcRoutes.HandleFunc("/number", s.logRequest(s.handleNumber("json")))
	//svcRoutes.HandleFunc("/number.txt", s.logRequest(s.handleNumber("txt")))
	//svcRoutes.HandleFunc("/alphanumeric", s.logRequest(s.handleAlphanumeric("json")))
	//svcRoutes.HandleFunc("/alphanumeric.txt", s.logRequest(s.handleAlphanumeric("txt")))
	//svcRoutes.HandleFunc("/ksuid", s.logRequest(s.handleKsuid("json")))
	//svcRoutes.HandleFunc("/ksuid.txt", s.logRequest(s.handleKsuid("txt")))
	//svcRoutes.HandleFunc("/uuid", s.logRequest(s.handleUuid("json")))
	//svcRoutes.HandleFunc("/uuid.txt", s.logRequest(s.handleUuid("txt")))
	//svcRoutes.HandleFunc("/dice", s.logRequest(s.handleDice("json")))
	//svcRoutes.HandleFunc("/dice.txt", s.logRequest(s.handleDice("txt")))
	//svcRoutes.HandleFunc("/coin", s.logRequest(s.handleCoin("json")))
	//svcRoutes.HandleFunc("/coin.txt", s.logRequest(s.handleCoin("txt")))
	//svcRoutes.Use(LogRequestMiddleware)

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
