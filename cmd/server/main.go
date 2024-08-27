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
	http.HandleFunc("/chore-list", handlers.ChoreListHandler(db))
	http.HandleFunc("/child-dashboard", handlers.ChildDashboardHandler(db))
	http.HandleFunc("/parent-panel", handlers.ParentPanelHandler(db))
	http.HandleFunc("/add-child", handlers.AddChildHandler(db))
	http.HandleFunc("/edit-child/{id}", handlers.EditChildHandler(db))
	http.HandleFunc("/delete-child/{id}", handlers.DeleteChildHandler(db))
	http.HandleFunc("/child-action", handlers.ChildActionHandler(db))
	http.HandleFunc("/add-chore", handlers.AddChoreHandler(db))
	http.HandleFunc("/edit-chore/{id}", handlers.EditChoreHandler(db))
	http.HandleFunc("/delete-chore/{id}", handlers.DeleteChoreHandler(db))
	http.HandleFunc("/chore-action", handlers.ChoreActionHandler(db))
	http.HandleFunc("/assign-chore", handlers.AssignChoreHandler(db))
	http.HandleFunc("/assign-chore-form", handlers.AssignChoreFormHandler(db))
	http.HandleFunc("/assignments-list", handlers.AssignmentsListHandler(db))
	http.HandleFunc("/assignment-action", handlers.AssignmentActionHandler(db))

	// Start the server
	fmt.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
