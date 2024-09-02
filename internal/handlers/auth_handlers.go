package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/slate20/goauth"
)

func LandingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("../../templates/public_layout.html", "../../templates/landing.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RegisterHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")

			err := auth.Register(username, email, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/public_layout.html", "../../templates/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func LoginHandler(db *sql.DB, auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			log.Println("Login attempt:", username)

			token, err := auth.Login(username, password)
			if err != nil {
				log.Println("Login failed:", err)
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			log.Println("Login successful")

			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    token,
				HttpOnly: true,
			})

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("../../templates/public_layout.html", "../../templates/login.html")
		if err != nil {
			log.Println("Login template error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Login template execution error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func ExtractUserID(r *http.Request, auth *goauth.AuthService) (int, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return 0, errors.New("Unauthorized")
	}

	token, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return int(userID), nil
}
