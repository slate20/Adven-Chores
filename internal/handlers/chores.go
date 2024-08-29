package handlers

import (
	"Adven-Chores/internal/models"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ChoreListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chores, err := models.GetAllChores(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/chore_list.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, chores)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AddChoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("../../templates/add_chore.html")
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
			description := r.FormValue("description")
			points, _ := strconv.Atoi(r.FormValue("points"))
			isRequired := r.FormValue("is_required") == "on"

			chore := &models.Chore{
				Description: description,
				Points:      points,
				IsRequired:  isRequired,
			}
			err := chore.Save(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Trigger", "refreshChoreList")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button class="action-button" hx-get="/add-chore" hx-target="#chore-action-container" hx-swap="innerHTML">Add Chore</button>`))
		}
	}
}

func EditChoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		id, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid chore ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			chore, err := models.GetChoreByID(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles("../../templates/edit_chore.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, chore)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			description := r.FormValue("description")
			points, _ := strconv.Atoi(r.FormValue("points"))
			isRequired := r.FormValue("is_required") == "on"

			chore := &models.Chore{
				ID:          id,
				Description: description,
				Points:      points,
				IsRequired:  isRequired,
			}
			err := chore.Save(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Trigger", "refreshChoreList")
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<button class="action-button" hx-get="/add-chore" hx-target="#chore-action-container" hx-swap="innerHTML">Add Chore</button>`))
		}
	}
}

func DeleteChoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		id, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid chore ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteChore(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ChoreActionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<button class="action-button" hx-get="/add-chore" hx-target="#chore-action-container" hx-swap="innerHTML">Add Chore</button>`))
	}
}

// Functions related to chore assignments

func AssignChoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		childID, err := strconv.ParseInt(r.FormValue("child_id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		choreID, err := strconv.ParseInt(r.FormValue("chore_id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid chore ID", http.StatusBadRequest)
			return
		}

		err = models.AssignChoreToChild(db, childID, choreID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Trigger refresh and return the Add Assignment button
		w.Header().Set("HX-Trigger", "refreshAssignments")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<button class="action-button" hx-get="/assign-chore-form" hx-target="#assignment-action-container" hx-swap="innerHTML">New Assignment</button>`))
	}
}

func AssignChoreFormHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Create a map of assigned chore IDs
		assignedChores := make(map[int64]bool)
		for _, assignment := range assignments {
			assignedChores[assignment.Chore.ID] = true
		}

		// Filter out assigned chores
		var availableChores []*models.Chore
		for _, chore := range chores {
			if !assignedChores[chore.ID] {
				availableChores = append(availableChores, chore)
			}
		}

		data := struct {
			Children []*models.Child
			Chores   []*models.Chore
		}{
			Children: children,
			Chores:   availableChores,
		}

		tmpl, err := template.ParseFiles("../../templates/assign_chore.html")
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

func AssignmentsListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assignments, err := models.GetAllAssignments(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/assignments_list.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, assignments)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AssignmentActionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<button class="action-button" hx-get="/assign-chore-form" hx-target="#assignment-action-container" hx-swap="innerHTML">New Assignment</button>`))
	}
}

func DeleteAssignmentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		assignmentID, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteAssignment(db, assignmentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CompleteAssignmentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		assignmentID, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
			return
		}

		err = models.CompleteAssignment(db, assignmentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Refresh the child-dashboard to update assignment status
		w.Header().Set("HX-Trigger", "refreshChildDashboard")
		w.Header().Set("Content-Type", "text/html")
	}
}

func RewardAssignmentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		assignmentID, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
			return
		}

		err = models.RewardAssignment(db, assignmentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Refresh the children list to update points
		w.Header().Set("HX-Trigger", "refreshList")
		w.Header().Set("Content-Type", "text/html")
	}
}

// function for child to accept a chore; takes in childID and choreID from the URL
func AcceptChoreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		childID, err := strconv.ParseInt(paths[len(paths)-2], 10, 64)
		if err != nil {
			http.Error(w, "Invalid child ID", http.StatusBadRequest)
			return
		}

		choreID, err := strconv.ParseInt(paths[len(paths)-1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid chore ID", http.StatusBadRequest)
			return
		}

		err = models.AssignChoreToChild(db, childID, choreID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Trigger refresh
		w.Header().Set("HX-Trigger", "refreshChildDashboard")
		w.Header().Set("Content-Type", "text/html")
	}
}
