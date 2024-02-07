package forum

import (
	"fmt"
	DB "forum/database"
	"net/http"
	"text/template"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/create_post.html"))
		tmpl.Execute(w, nil)
		return
	}

	title := r.FormValue("title")
	category := r.FormValue("category")
	content := r.FormValue("content")

	var db = DB.OpenDataBase()
	defer db.Close()

	_, err := db.Exec("INSERT INTO posts (title, content, category, likes) VALUES (?, ?, ?, 0)", title, content, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?category=%s", category), http.StatusFound)
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "La catégorie est requise.", http.StatusBadRequest)
		return
	}

	posts, err := FetchPostsByCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmpl *template.Template

	switch category {
	case "Général":
		tmpl = template.Must(template.ParseFiles("./templates/categories/general.html"))
	case "Informatique":
		tmpl = template.Must(template.ParseFiles("./templates/categories/informatique.html"))
	case "Jeux Videos":
		tmpl = template.Must(template.ParseFiles("./templates/categories/jeux_videos.html"))
	case "Musique":
		tmpl = template.Must(template.ParseFiles("./templates/categories/musique.html"))
	case "Actualité":
		tmpl = template.Must(template.ParseFiles("./templates/categories/actualite.html"))
	case "Evenement de l'année":
		tmpl = template.Must(template.ParseFiles("./templates/categories/evenement_de_l_annee.html"))
	case "Technologie":
		tmpl = template.Must(template.ParseFiles("./templates/categories/technologie.html"))
	case "Art et culture":
		tmpl = template.Must(template.ParseFiles("./templates/categories/art_et_culture.html"))
	case "Cinema":
		tmpl = template.Must(template.ParseFiles("./templates/categories/cinema.html"))
	default:
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	tmpl.Execute(w, struct {
		Category string
		Posts    []Post
	}{
		Category: category,
		Posts:    posts,
	})
}

func FetchPostByID(postID int) (Post, error) {
	var db = DB.OpenDataBase()
	defer db.Close()

	var post Post
	err := db.QueryRow("SELECT id, title, content, category, likes FROM posts WHERE id = ?", postID).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func FetchPostsByCategory(category string) ([]Post, error) {
	var db = DB.OpenDataBase()
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content, category, likes FROM posts WHERE category = ? ORDER BY id DESC", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
