package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const VERSION = "0.0.1"

func main() {
	startTime := time.Now()

	flagVersion := flag.Bool("v", false, "version")
	flagDebug := flag.Bool("debug", false, "Printing all log levels")
	flagDev := flag.Bool("dev", false, "loggs include caller and debug level")
	flag.Parse()

	if *flagVersion {
		fmt.Println("Version: " + VERSION)
		os.Exit(0)
	}


	logger := log.New(os.Stderr)
	logger.SetPrefix("main")
	logger.SetReportTimestamp(true)

	if *flagDebug {
		logger.SetLevel(log.DebugLevel)
	}
	if *flagDev {
		logger.SetReportCaller(true)
	}

	l := logger.WithPrefix(logger.GetPrefix() + ".http")
	router := newRouter(l)
	router.prepareRoutes()

	logger.Info(
		"handing over to http server",
		logTimeSpent, time.Since(startTime),
	)

	logger.Fatal(l, runServer(l, &router))
}

func runServer(l *log.Logger, s *server) error {

	pageSrv := http.Server{
		Addr:    ":8080",
		Handler: s.Handler(),
	}

	l.Info("Starting server", "addr", pageSrv.Addr)
	if err := pageSrv.ListenAndServe(); err != nil {
		return fmt.Errorf("ListenAndServe: %s", err)
	}
	return fmt.Errorf("unexpected server shutdown")
}
