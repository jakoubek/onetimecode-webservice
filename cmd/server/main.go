package main

import (
	"expvar"
	"flag"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type config struct {
	port        int
	env         string
	logfilePath string
	limiter     struct {
		rps     float64
		burst   int
		enabled bool
	}
	securekey string
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

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 1, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 3, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.securekey, "securekey", "", "Securekey for accessing the metrics endpoint")

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

	expvar.NewString("version").Set(version)
	expvar.NewString("buildVersion").Set(buildVersion)
	expvar.NewString("buildTime").Set(buildTime)
	expvar.NewString("serverStartupTime").Set(app.startupTime.String())
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	err = app.serve()
	if err != nil {
		logger.Fatalln(err)
	}

}
