package handlers

import (
	"Adven-Chores/internal/models"
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"github.com/slate20/goauth"
)

func HomeHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
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

		var firstChildID string
		if len(children) > 0 {
			firstChildID = strconv.FormatInt(children[0].ID, 10)
		}

		data := struct {
			HasChildren  bool
			FirstChildID string
		}{
			HasChildren:  len(children) > 0,
			FirstChildID: firstChildID,
		}

		tmpl, err := template.ParseFiles("../../templates/layout.html")
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
