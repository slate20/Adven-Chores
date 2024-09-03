package handlers

import (
	"Adven-Chores/internal/models"
	"database/sql"
	"html/template"
	"net/http"

	"github.com/slate20/goauth"
)

func ChildNavHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := ExtractUserID(r, auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		children, err := models.GetChildrenByUserID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/child_nav.html")
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
