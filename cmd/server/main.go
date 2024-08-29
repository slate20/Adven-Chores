package main

import (
	"Adven-Chores/internal/database"
	"Adven-Chores/internal/handlers"
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

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	// Webserver routes
	http.HandleFunc("/", handlers.HomeHandler(db))
	http.HandleFunc("/child-nav", handlers.ChildNavHandler(db))
	http.HandleFunc("/child-list", handlers.ChildListHandler(db))
	http.HandleFunc("/chore-list", handlers.ChoreListHandler(db))
	http.HandleFunc("/child-dashboard/{id}", handlers.ChildDashboardHandler(db))
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
	http.HandleFunc("/delete-assignment/{id}", handlers.DeleteAssignmentHandler(db))
	http.HandleFunc("/assignment-action", handlers.AssignmentActionHandler(db))
	http.HandleFunc("/accept-chore/{child_id}/{chore_id}", handlers.AcceptChoreHandler(db))
	http.HandleFunc("/complete-chore/{id}", handlers.CompleteAssignmentHandler(db))
	http.HandleFunc("/reward-assignment/{id}", handlers.RewardAssignmentHandler(db))
	http.HandleFunc("/reward-list", handlers.RewardListHandler(db))
	http.HandleFunc("/add-reward", handlers.AddRewardHandler(db))
	http.HandleFunc("/edit-reward/{id}", handlers.EditRewardHandler(db))
	http.HandleFunc("/delete-reward/{id}", handlers.DeleteRewardHandler(db))
	http.HandleFunc("/reward-action", handlers.RewardActionHandler(db))
	http.HandleFunc("/rewards-store/{child_id}", handlers.RewardsStoreHandler(db))
	http.HandleFunc("/redeem-reward/", handlers.RedeemRewardHandler(db))
	http.HandleFunc("/set-pin", handlers.SetPinHandler(db))

	// Start the server
	fmt.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
