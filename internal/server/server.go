package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	repo "jst.dev/internal/persistence"
)

type Server struct {
	port int

	db *repo.TursoRepo
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to parse PORT: %v", err)
	}
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatalf("HOST is not set")
	}

	db, err := repo.NewTursoRepo()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	NewServer := &Server{
		port: port,
		db:   db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
