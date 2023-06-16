package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/matryer/way"
)

func main() {
	flagPort := flag.Int("port", 8000, "port to serve on")
	flag.Parse()

	server := &server{
		router: way.NewRouter(),
	}

	server.routes()

	listenAddr := fmt.Sprintf(":%d", *flagPort)
	log.Printf("[Reverse Proxy]: Listening on %s...\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, server.router))
}

type server struct {
	router *way.Router
}

// Register handlers for routes
func (srv *server) routes() {
	srv.router.Handle("GET", "git.jst.dev/", srv.gitHandler())

}

func (srv *server) gitHandler() *httputil.ReverseProxy {
	// setup
	urlPortfolio, err := url.Parse("https://jst.dev")
	if err != nil {
		log.Fatal(err)
	}

	// handler
	return newProxy(urlPortfolio)
}

func newProxy(target *url.URL) *httputil.ReverseProxy {

	errorLog := log.New(log.Writer(), "proxy error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.RequestURI = ""

		// To prevent IP spoofing, be sure to delete any pre-existing
		// X-Forwarded-For header coming from the client or
		// an untrusted proxy.
		req.Header.Del("X-Forwarded-For")
	}

	return &httputil.ReverseProxy{
		Director:      director,
		FlushInterval: 50 * time.Millisecond,
		ErrorLog:      errorLog,
	}
}
