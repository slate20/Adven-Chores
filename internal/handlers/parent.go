package handlers

import (
	"Adven-Chores/internal/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/slate20/goauth"
)

func ParentPanelHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := ExtractUserID(r, auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// get parent_pin from database users table
		var parentPin int
		err = db.QueryRow("SELECT parent_pin FROM users WHERE id = ?", userID).Scan(&parentPin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// convert hx-prompt header to int and check if it matches parent_pin
		pin, _ := strconv.Atoi(r.Header.Get("Hx-Prompt"))
		if pin != parentPin {
			http.Error(w, "Incorrect PIN", http.StatusUnauthorized)
			return
		}

		children, err := models.GetChildrenByUserID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		chores, err := models.GetChoresByUserID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		assignments, err := models.GetAssignmentsByUserID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rewards, err := models.GetRewardsByUserID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Children    []*models.Child
			Chores      []*models.Chore
			Assignments []*models.Assignment
			Rewards     []*models.Reward
		}{
			Children:    children,
			Chores:      chores,
			Assignments: assignments,
			Rewards:     rewards,
		}

		tmpl, err := template.ParseFiles(
			"../../templates/parent_panel.html",
			"../../templates/child_list.html",
			"../../templates/chore_list.html",
			"../../templates/assignments_list.html",
			"../../templates/reward_list.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "parent_panel.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// function to handle set-pin form and set parentPin
func SetPinHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := ExtractUserID(r, auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("SetPinHandler: method=%s", r.Method)
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("../../templates/set_pin.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			// convert form input to int and save to database
			pin, _ := strconv.Atoi(r.FormValue("pin"))
			_, err = db.Exec("UPDATE users SET parent_pin = ? WHERE id = ?", pin, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Return the set-pin button
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button id="set-pin-btn" hx-get="/set-pin" hx-target="#set-pin-container" hx-swap="innerHTML">Change Pin</button>`))
			log.Printf("SetPinHandler: new pin set to %d", pin)
		}
	}
}

// Handler for account settings page
func AccountSettingsHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := ExtractUserID(r, auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user := &models.User{}
		err = db.QueryRow("SELECT id, username, email, parent_pin FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email, &user.ParentPin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/account_settings.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
