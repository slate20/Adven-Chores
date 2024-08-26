package main

import (
	"ChoreQuest/internal/database"
	"ChoreQuest/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	// TODO: Connect to the database
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Webserver routes
	http.HandleFunc("/", handlers.HomeHandler(db))
	http.HandleFunc("/child-list", handlers.ChildListHandler(db))
	http.HandleFunc("/child-dashboard", handlers.ChildDashboardHandler(db))
	http.HandleFunc("/parent-panel", handlers.ParentPanelHandler(db))
	http.HandleFunc("/add-child", handlers.AddChildHandler(db))
	http.HandleFunc("/edit-child/{id}", handlers.EditChildHandler(db))
	http.HandleFunc("/delete-child/{id}", handlers.DeleteChildHandler(db))
	http.HandleFunc("/child-action", handlers.ChildActionHandler(db))

	// Start the server
	fmt.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
