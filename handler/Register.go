package forum

import (
	"database/sql"
	DB "forum/database"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	var db = DB.OpenDataBase()
	defer db.Close()

	var existingUser string
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUser)

	if err != sql.ErrNoRows {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		tmpl.Execute(w, struct {
			Error string
		}{
			Error: "Le pseudo que vous avez choisi est déjà utilisé.",
		})
		return
	}

	var existingEmail string
	errEmail := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)

	if errEmail != sql.ErrNoRows {
		tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		tmpl.Execute(w, struct {
			Error string
		}{
			Error: "Cet email est déjà utilisé.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
