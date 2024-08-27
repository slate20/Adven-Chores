package handlers

import (
	"ChoreQuest/internal/models"
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
			w.Write([]byte(`<button hx-get="/add-chore" hx-target="#chore-action-container" hx-swap="innerHTML">Add Chore</button>`))
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
			w.Write([]byte(`<button hx-get="/add-chore" hx-target="#chore-action-container" hx-swap="innerHTML">Add Chore</button>`))
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
		w.Write([]byte(`<button hx-get="/add-chore" hx-target="#chore-action-container" hx-swap="innerHTML">Add Chore</button>`))
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
		w.Write([]byte(`<button hx-get="/assign-chore-form" hx-target="#assignments-form" hx-swap="innerHTML">Assign Chore</button>`))
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

		data := struct {
			Children []*models.Child
			Chores   []*models.Chore
		}{
			Children: children,
			Chores:   chores,
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
		w.Write([]byte(`<button hx-get="/assign-chore-form" hx-target="#assignment-action-container" hx-swap="innerHTML">New Assignment</button>`))
	}
}
