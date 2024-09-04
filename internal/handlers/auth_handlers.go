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

func LandingHandler(auth *goauth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			_, err = auth.ValidateToken(cookie.Value)
			if err == nil {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
		}

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
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			_, err = auth.ValidateToken(cookie.Value)
			if err == nil {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm-password")

			// Check if passwords match
			if password != confirmPassword {
				if r.Header.Get("HX-Request") == "true" {
					w.Write([]byte("<div id='error-message'>Passwords do not match</div>"))
				} else {
					http.Error(w, "Passwords do not match", http.StatusBadRequest)
				}
				return
			}

			// Check if user already exists
			var userexists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)", username, email).Scan(&userexists)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if userexists {
				if r.Header.Get("HX-Request") == "true" {
					w.Write([]byte("<div id='error-message'>A user already exists with that username or email</div>"))
				} else {
					http.Error(w, "A user already exists with that username or email", http.StatusBadRequest)
				}
				return
			}

			err := auth.Register(username, email, password)
			if err != nil {
				if r.Header.Get("HX-Request") == "true" {
					w.Write([]byte("<div class='error-message'>" + err.Error() + "</div>"))
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				return
			}

			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
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
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			_, err = auth.ValidateToken(cookie.Value)
			if err == nil {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			log.Println("Login attempt:", username)

			token, err := auth.Login(username, password)
			if err != nil {
				log.Println("Login failed:", err)
				if r.Header.Get("HX-Request") == "true" {
					w.Write([]byte("<div id='error-message'>Invalid username or password</div>"))
				} else {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				}
				return
			}

			log.Println("Login successful")

			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    token,
				HttpOnly: true,
			})

			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/home")
			} else {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}
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
