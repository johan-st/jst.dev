package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
)

func routes(l *log.Logger) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome(l))
	mux.HandleFunc("/favicon.ico", handleFile(l, "./content/public/favicon.ico"))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./content/public")))) //TODO: do not list files

	mux.HandleFunc("/404", handleNotFound(l))

	return mux
}

func handleHome(l *log.Logger) http.HandlerFunc {
	notFound := handleNotFound(l)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFound(w, r)
			return
		}

		l.Println("Serving home")
		html, err := os.ReadFile("./content/html/index.html")
		if err != nil {
			l.Fatal("handleHome: ", err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
	}
}

func handleNotFound(l *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Printf("Serving 404 on path: %s. Referer: %s\n", r.URL.Path, r.Referer())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found\nReferer: " + r.Referer()))
	}
}

// GENRALIZED HANDLERS

func handleFile(l *log.Logger, filename string) http.HandlerFunc {
	bytes, mimeType, err := fileBytesMime(filename)
	if err != nil {
		l.Fatal("handleFile: ", err)
	}

	gzipBytes, err := gzipBytes(bytes)
	if err != nil {
		l.Fatal("handleFile: ", err)
	}

	l.Printf("Serving file: %s,  mime: %s, size: %d, gzipped: %d\n", filename, mimeType, len(bytes), len(gzipBytes))

	return func(w http.ResponseWriter, r *http.Request) {
		l.Println("Serving file:", filename)
		w.Header().Set("Content-Type", mimeType)
		if r.Header.Get("Accept-Encoding") == "gzip" {
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(gzipBytes)
			return
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
		w.Write(bytes)
	}
}

// HELPER FUNCTIONS

func fileBytesMime(filename string) ([]byte, string, error) {
	const MaxInMemSize = 1024 * 1024 // 1MB
	var (
		bytes    []byte
		mimeType string
		err      error
	)

	file, err := os.Open(filename)
	if err != nil {
		return bytes, mimeType, fmt.Errorf("error opening file (%s): %s", filename, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return bytes, mimeType, fmt.Errorf("error stat file (%s): %s", filename, err)
	}

	if stat.IsDir() {
		return bytes, mimeType, fmt.Errorf("error opening file (%s): is a directory", filename)
	}

	if stat.Size() == 0 {
		return bytes, mimeType, fmt.Errorf("error opening file (%s): file is empty", filename)
	}

	if stat.Size() > MaxInMemSize {
		return bytes, mimeType, fmt.Errorf("error opening file (%s): file is too large. Max allowed (%d B)", filename, MaxInMemSize)
	}

	bytes = make([]byte, stat.Size())
	_, err = file.Read(bytes)
	if err != nil {
		return bytes, mimeType, fmt.Errorf("error reading file (%s): %e", filename, err)
	}

	mimeType = mime.TypeByExtension(filename)
	if mimeType == "" {
		mimeType = http.DetectContentType(bytes)
	}
	if mimeType == "" {
		return bytes, mimeType, fmt.Errorf("error reading file (%s): unknown file type ", filename)
	}

	return bytes, mimeType, nil
}

func gzipBytes(bs []byte) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
		zw  *gzip.Writer
	)

	zw, err = gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return bs, err
	}

	_, err = zw.Write(bs)
	if err != nil {
		return bs, err
	}

	if err := zw.Close(); err != nil {
		return bs, err
	}

	return buf.Bytes(), nil
}
