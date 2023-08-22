package main

import (
	"embed"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/matryer/way"
)

//go:embed template assets
var staticFS embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	flagDebug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(true)
	if *flagDebug {
		logger.SetLevel(log.DebugLevel)
	}

	handler := handler{
		l:      logger,
		fs:     staticFS,
		router: way.NewRouter(),
	}

	handler.routes()

	httpServer := http.Server{
		Addr:         ":8080",
		Handler:      &handler,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Infof("Listening on %s", httpServer.Addr)
	return httpServer.ListenAndServe()
}
