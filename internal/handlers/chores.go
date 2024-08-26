package handlers

import (
	"ChoreQuest/internal/models"
	"database/sql"
	"html/template"
	"net/http"
)

func ChoresHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chores, err := models.GetAllChores(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/chores.html")
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
