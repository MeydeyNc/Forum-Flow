package forum

import (
	DB "forum/database"
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.ServeFile(w, r, "./templates/noconnected.html")
		return
	}

	var db = DB.OpenDataBase()
	defer db.Close()

	rows, err := db.Query("SELECT title, content, category FROM posts ORDER BY id DESC LIMIT 3")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Title, &post.Content, &post.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	tmpl := template.Must(template.ParseFiles("./templates/home.html"))
	tmpl.Execute(w, struct {
		Username string
		Posts    []Post
	}{
		Username: cookie.Value,
		Posts:    posts,
	})
}
