package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

func main() {

	logger := log.New(os.Stderr)
	logger.SetPrefix("main")
	logger.SetReportTimestamp(true)
	logger.SetLevel(log.DebugLevel)
	logger.SetReportCaller(true)

	router := newRouter(logger)
	router.prepareRoutes()

	log.Fatal(runServer(&router))
}

func runServer(s *server) error {
	l := s.l.WithPrefix("http")

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
