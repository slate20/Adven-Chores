package handlers

import (
	"Adven-Chores/internal/models"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ChildListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		children, err := models.GetAllChildren(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/child_list.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, children)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ChildDashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/child-dashboard/"):]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		child, err := models.GetChildByID(db, id)
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

		childAssignments, err := models.GetAssignmentsByChild(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// create a map of assigned chore IDs and filter to get only unassigned chores
		assignedChores := make(map[int64]bool)
		for _, assignment := range assignments {
			assignedChores[assignment.Chore.ID] = true
		}

		var availableChores []*models.Chore
		for _, chore := range chores {
			if !assignedChores[chore.ID] {
				availableChores = append(availableChores, chore)
			}
		}

		data := struct {
			Child           *models.Child
			Assignments     []*models.Assignment
			AvailableChores []*models.Chore
		}{
			Child:           child,
			Assignments:     childAssignments,
			AvailableChores: availableChores,
		}

		tmpl, err := template.ParseFiles("../../templates/child_dashboard.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Function to add a new child
func AddChildHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("../../templates/add_child.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
		} else if r.Method == http.MethodPost {
			name := r.FormValue("name")
			child := &models.Child{Name: name, Points: 0}
			err := child.Save(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Trigger refresh and return the Add Child button
			w.Header().Set("HX-Trigger", "refreshList")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button class="action-button" hx-get="/add-child" hx-target="#child-action-container" hx-swap="innerHTML">Add Child</button>`))
		}
	}
}

// Function to load exisiting child data for editing
func EditChildHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		id, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			child, err := models.GetChildByID(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles("../../templates/edit_child.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, child)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			name := r.FormValue("name")
			points, err := strconv.Atoi(r.FormValue("points"))
			rewards := r.FormValue("rewards")
			if err != nil {
				http.Error(w, "Invalid points value", http.StatusBadRequest)
				return
			}

			child := &models.Child{ID: id, Name: name, Points: points, Rewards: rewards}
			err = child.Save(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Trigger refresh and return the Add Child button
			w.Header().Set("HX-Trigger", "refreshList")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button class="action-button" hx-get="/add-child" hx-target="#child-action-container" hx-swap="innerHTML">Add Child</button>`))
		}
	}
}

// Function to delete a child
func DeleteChildHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		id, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteChild(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		// Trigger refresh
		w.Header().Set("HX-Trigger", "refreshList")
		w.Header().Set("Content-Type", "text/html")
	}
}

func ChildActionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<button class="action-button" hx-get="/add-child" hx-target="#child-action-container" hx-swap="innerHTML">Add Child</button>`))
	}
}
