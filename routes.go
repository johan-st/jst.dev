package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	// config Config,
	// tenantsStore        *TenantsStore,
	// commentsStore       *CommentsStore,
	// conversationService *ConversationService,
	// chatGPTService      *ChatGPTService,
	// authProxy           *authProxy
) {
	mux.Handle("/hello", handleHello(logger))
	mux.Handle("/maybe", onlySometimes(handleAdmin(logger)))
	mux.Handle("/", http.NotFoundHandler())
}

func handleHello(logger *log.Logger) http.Handler {
	// closure for setting up the handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("hello")
		w.Write([]byte("hello"))
	})
}

func handleAdmin(logger *log.Logger) http.Handler {
	// closure for setting up the handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("admin")
		encode(w, r, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "admin page was hit"})
	})
}

// MIDDLWARE

func onlySometimes(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAdmin {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// HELPER FUNCTIONS

func currentUser(_ *http.Request) User {
	// get the user from the request
	if rand.Intn(2) == 0 {
		return User{IsAdmin: true}
	}
	return User{IsAdmin: false}
}

type User struct {
	IsAdmin bool
}

// err := encode(w, r, http.StatusOK, obj)
func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// Since it is a return argument in decode, you will need to specify the type you expect:
// go decoded, err := decode[CreateSomethingRequest](r)
func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// // Validator is an object that can be validated.
// type Validator interface {
// 	// Valid checks the object and returns any
// 	// problems. If len(problems) == 0 then
// 	// the object is valid.
// 	Valid(ctx context.Context) (problems map[string]string)
// }
// func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
// 	var v T
// 	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
// 		return v, nil, fmt.Errorf("decode json: %w", err)
// 	}
// 	if problems := v.Valid(r.Context()); len(problems) > 0 {
// 		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
// 	}
// 	return v, nil, nil
// }
