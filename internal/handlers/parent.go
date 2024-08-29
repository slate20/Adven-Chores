package handlers

import (
	"Adven-Chores/internal/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var parentPin = "1234"

func ParentPanelHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for correct pin in Hx-Prompt header
		pin := r.Header.Get("Hx-Prompt")
		if pin != parentPin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		children, err := models.GetAllChildren(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		chores, err := models.GetAllChores(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		assignments, err := models.GetAllAssignments(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rewards, err := models.GetAllRewards(db)
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
func SetPinHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			pin := r.FormValue("pin")
			parentPin = pin

			// Return the set-pin button
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button id="set-pin-btn" hx-get="/set-pin" hx-target="#set-pin-container" hx-swap="innerHTML">Change Pin</button>`))
			log.Printf("SetPinHandler: new pin set to %s", pin)
		}
	}
}
