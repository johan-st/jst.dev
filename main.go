package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	startTime := time.Now()

	flagVerbose := flag.Bool("v", false, "Printing all log levels")
	flagDev := flag.Bool("dev", false, "loggs include caller")
	flag.Parse()

	logger := log.New(os.Stderr)
	logger.SetPrefix("main")
	logger.SetReportTimestamp(true)

	if *flagVerbose {
		logger.SetLevel(log.DebugLevel)
	}
	if *flagDev {
		logger.SetReportCaller(true)
	}
	
	l := logger.WithPrefix(logger.GetPrefix() + ".http")
	router := newRouter(l)
	router.prepareRoutes()

	logger.Info(
		"startup complete, handing over to http server",
		"time", time.Since(startTime),
	)
	
	logger.Fatal(l, runServer(l,&router))
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
