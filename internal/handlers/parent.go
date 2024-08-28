package handlers

import (
	"ChoreQuest/internal/models"
	"database/sql"
	"html/template"
	"net/http"
)

func ParentPanelHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for correct pin in Hx-Prompt header
		pin := r.Header.Get("Hx-Prompt")
		if pin != "9228" {
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
