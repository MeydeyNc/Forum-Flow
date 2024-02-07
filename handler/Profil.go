package forum

import (
	"net/http"
	"text/template"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, email, numPosts, numLikes, numDislikes := GetUserData(cookie.Value)

	data := struct {
		Username    string
		Email       string
		NumPosts    int
		NumLikes    int
		NumDislikes int
	}{
		Username:    username,
		Email:       email,
		NumPosts:    numPosts,
		NumLikes:    numLikes,
		NumDislikes: numDislikes,
	}

	tmpl, _ := template.ParseFiles("./templates/profil.html")
	tmpl.Execute(w, data)
}
