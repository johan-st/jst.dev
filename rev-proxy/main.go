package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	// env vars
	useTls := os.Getenv("TLS_ENABLED")
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")
	port := os.Getenv("PORT")

	// default values
	if tlsCert == "" {
		tlsCert = "cert.pem"
	}
	if tlsKey == "" {
		tlsKey = "key.pem"
	}
	if port == "" {
		port = "8080"
	}

	// server
	// server := &server{
	// 	router: way.NewRouter(),
	// }

	// server.routes()

	// routes
	mux := http.NewServeMux()
	mux.Handle("git.jst.dev/", gitHandler())
	mux.Handle("/git", gitHandler())
	mux.Handle("me.jst.dev/", portfolioHandler())
	mux.Handle("/me", portfolioHandler())
	mux.Handle("/", notFoundHandler())

	listenAddr := fmt.Sprintf(":%s", port)
	log.Printf("[Reverse Proxy]: Listening on %s...\n", listenAddr)
	if useTls != "" && useTls != "false" && useTls != "0" && useTls != "no" {
		log.Println("Using TLS")
		log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCert, tlsKey, mux))
	} else {
		log.Println("Not using TLS")
		log.Fatal(http.ListenAndServe(listenAddr, mux))
	}
}

func portfolioHandler() *httputil.ReverseProxy {
	// setup
	urlPortfolio, err := url.Parse("https://jst.dev/")
	if err != nil {
		log.Fatal(err)
	}

	// handler
	return newProxy(urlPortfolio)
}

func gitHandler() *httputil.ReverseProxy {
	// setup
	urlPortfolio, err := url.Parse("https://github.com/")
	if err != nil {
		log.Fatal(err)
	}

	// handler
	return newProxy(urlPortfolio)
}

func notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
		log.Printf("Not Found: %s\n", r.URL.Path)
	}
}

func newProxy(target *url.URL) *httputil.ReverseProxy {

	errorLog := log.New(log.Writer(), "proxy error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	director := func(req *http.Request) {
		log.Printf("Proxying %s to %s\n", req.URL.String(), target.String()+req.URL.Path)
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
