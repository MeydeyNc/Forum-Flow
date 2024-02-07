package forum

import (
	DB "forum/database"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var db = DB.OpenDataBase()
	defer db.Close()

	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Eh non": "Mauvais pseudo ou mot de passe",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Eh non": "Mauvais pseudo ou mot de passe",
		})

		return
	}

	usernameCookie := &http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	}
	http.SetCookie(w, usernameCookie)

	http.Redirect(w, r, "/home", http.StatusFound)
}
