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

const (
	URLGit       = "https://github.com"
	URLPortfolio = "https://jst.dev/"
	URLImg       = "https://jst.dev:8080"

	DefaultCert    = "cert.pem"
	DefaultKey     = "key.pem"
	DefaultPort    = "8080"
	DefaultLogfile = "proxy.log"
)

func main() {
	// env vars
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")
	useTls := os.Getenv("PROXY_TLS_ENABLED")
	port := os.Getenv("PROXY_PORT")
	pathLogfile := os.Getenv("PROXY_pathLOGFILE")

	// default values
	if tlsCert == "" {
		tlsCert = DefaultCert
	}
	if tlsKey == "" {
		tlsKey = DefaultKey
	}
	if port == "" {
		port = DefaultPort
	}
	if pathLogfile == "" {
		pathLogfile = DefaultLogfile
	}

	// logging
	logfile, err := os.OpenFile(DefaultLogfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	// routes
	mux := http.NewServeMux()
	mux.Handle("git.jst.dev/", handlerGit())
	mux.Handle("/git/", handlerGit())

	mux.Handle("me.jst.dev/", handlerPortfolio())
	mux.Handle("/me/", handlerPortfolio())

	mux.Handle("img.jst.dev/", handlerImg())
	mux.Handle("/img/", handlerImg())

	mux.Handle("/", handlerNotFound())

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

// HANDLERS

func handlerPortfolio() *httputil.ReverseProxy {
	// setup
	urlProxy, err := url.Parse(URLPortfolio)
	if err != nil {
		log.Fatal(err)
	}

	// handler
	return newReverserProxy(urlProxy)
}

func handlerGit() *httputil.ReverseProxy {
	// setup
	urlProxy, err := url.Parse(URLGit)
	if err != nil {
		log.Fatal(err)
	}

	// handler
	return newReverserProxy(urlProxy)
}

func handlerImg() *httputil.ReverseProxy {
	// setup
	urlProxy, err := url.Parse(URLImg)
	if err != nil {
		log.Fatal(err)
	}

	// handler
	return newReverserProxy(urlProxy)
}

func handlerNotFound() http.HandlerFunc {
	// setup

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
		log.Printf("Not Found: %s\n", r.URL.Path)
	}
}

// PROXY
func newReverserProxy(target *url.URL) *httputil.ReverseProxy {

	errorLog := log.New(log.Writer(), "proxy error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	director := func(req *http.Request) {
		log.Printf("Proxying:\t%s  ->  %s\n", req.URL.String(), target.String())

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
