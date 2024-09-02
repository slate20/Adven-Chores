package main

import (
	"Adven-Chores/internal/database"
	"Adven-Chores/internal/handlers"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/slate20/goauth"
)

var auth *goauth.AuthService

func main() {
	// Connect to the database
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Initialize GoAuth
	jwtSecret := os.Getenv("JWT_SECRET")
	auth, err = goauth.NewAuthService(db, jwtSecret, 24*time.Hour, 1*time.Hour)
	if err != nil {
		fmt.Println("Error initializing GoAuth:", err)
		return
	}

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	// Public routes
	http.HandleFunc("/", handlers.LandingHandler())
	http.HandleFunc("/register", handlers.RegisterHandler(db, auth))
	http.HandleFunc("/login", handlers.LoginHandler(db, auth))
	http.HandleFunc("/logout", handlers.LogoutHandler())

	// Protected routes
	http.HandleFunc("/home", authMiddleware(handlers.HomeHandler(db)))
	http.HandleFunc("/child-nav", authMiddleware(handlers.ChildNavHandler(db)))
	http.HandleFunc("/child-list", authMiddleware(handlers.ChildListHandler(db)))
	http.HandleFunc("/chore-list", authMiddleware(handlers.ChoreListHandler(db)))
	http.HandleFunc("/child-dashboard/{id}", authMiddleware(handlers.ChildDashboardHandler(db)))
	http.HandleFunc("/parent-panel", authMiddleware(handlers.ParentPanelHandler(db)))
	http.HandleFunc("/add-child", authMiddleware(handlers.AddChildHandler(db)))
	http.HandleFunc("/edit-child/{id}", authMiddleware(handlers.EditChildHandler(db)))
	http.HandleFunc("/delete-child/{id}", authMiddleware(handlers.DeleteChildHandler(db)))
	http.HandleFunc("/child-action", authMiddleware(handlers.ChildActionHandler(db)))
	http.HandleFunc("/add-chore", authMiddleware(handlers.AddChoreHandler(db)))
	http.HandleFunc("/edit-chore/{id}", authMiddleware(handlers.EditChoreHandler(db)))
	http.HandleFunc("/delete-chore/{id}", authMiddleware(handlers.DeleteChoreHandler(db)))
	http.HandleFunc("/chore-action", authMiddleware(handlers.ChoreActionHandler(db)))
	http.HandleFunc("/assign-chore", authMiddleware(handlers.AssignChoreHandler(db)))
	http.HandleFunc("/assign-chore-form", authMiddleware(handlers.AssignChoreFormHandler(db)))
	http.HandleFunc("/assignments-list", authMiddleware(handlers.AssignmentsListHandler(db)))
	http.HandleFunc("/delete-assignment/{id}", authMiddleware(handlers.DeleteAssignmentHandler(db)))
	http.HandleFunc("/assignment-action", authMiddleware(handlers.AssignmentActionHandler(db)))
	http.HandleFunc("/accept-chore/{child_id}/{chore_id}", authMiddleware(handlers.AcceptChoreHandler(db)))
	http.HandleFunc("/complete-chore/{id}", authMiddleware(handlers.CompleteAssignmentHandler(db)))
	http.HandleFunc("/reward-assignment/{id}", authMiddleware(handlers.RewardAssignmentHandler(db)))
	http.HandleFunc("/reward-list", authMiddleware(handlers.RewardListHandler(db)))
	http.HandleFunc("/add-reward", authMiddleware(handlers.AddRewardHandler(db)))
	http.HandleFunc("/edit-reward/{id}", authMiddleware(handlers.EditRewardHandler(db)))
	http.HandleFunc("/delete-reward/{id}", authMiddleware(handlers.DeleteRewardHandler(db)))
	http.HandleFunc("/reward-action", authMiddleware(handlers.RewardActionHandler(db)))
	http.HandleFunc("/rewards-store/{child_id}", authMiddleware(handlers.RewardsStoreHandler(db)))
	http.HandleFunc("/redeem-reward/", authMiddleware(handlers.RedeemRewardHandler(db)))
	http.HandleFunc("/set-pin", authMiddleware(handlers.SetPinHandler(db)))

	// Start the server
	fmt.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token, err := auth.ValidateToken(cookie.Value)
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
