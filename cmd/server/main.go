package main

import (
	"ChoreQuest/internal/handlers"
	"fmt"
	"net/http"
)

type Chore struct {
	Name  string
	Type  string
	Value int
}

type Child struct {
	Name   string
	Job    string
	Points int
}

type Reward struct {
	Name  string
	Value int
}

func main() {
	// TODO: Connect to the database

	// Webserver routes
	http.HandleFunc("/", handlers.HomeHandler)

	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
