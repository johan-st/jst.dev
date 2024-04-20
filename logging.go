package main

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

type EventFrontend struct {
	Type string `json:"type"`
	User string `json:"user"`
	Path string `json:"path"`
	Message string `json:"message"`
}


func handlerLogging(e echo.Context) error {
	var event EventFrontend
	fmt.Println("Logging request")
	e.Bind(&event)

	fmt.Printf("Type: %s\n", event.Type)
	fmt.Printf("User: %s\n", event.User)
	fmt.Printf("Path: %s\n", event.Path)
	fmt.Printf("Message: %s\n\n", event.Message)

	
    return nil
}