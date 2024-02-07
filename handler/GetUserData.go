package forum

import (
	DB "forum/database"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func GetUserData(username string) (string, string, int, int, int) {
	var db = DB.OpenDataBase()
	defer db.Close()

	var user User
	err := db.QueryRow("SELECT id, username, email, password, numPosts, numLikes, numDislikes FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.NumPosts, &user.NumLikes, &user.NumDislikes)
	if err != nil {
		log.Fatal(err)
	}

	return user.Username, user.Email, user.NumPosts, user.NumLikes, user.NumDislikes
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/user/")

	_, email, numPosts, numLikes, numDislikes := GetUserData(username)

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
