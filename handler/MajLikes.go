package forum

import (
	"fmt"
	DB "forum/database"
	"net/http"
	"strconv"
)

func MajLikesHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.FormValue("postid")
	if len(postIDStr) == 0 {
		http.Error(w, "Manque id post", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Id post invalide", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")

	var db = DB.OpenDataBase()
	defer db.Close()

	post, err := FetchPostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "Veuillez vous connecter.", http.StatusUnauthorized)
		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND username = ?", postID, cookie.Value).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if action == "like" {
		if count > 0 {
			http.Error(w, "Vous avez déjà liké ce post", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO likes (post_id, username) VALUES (?, ?)", postID, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE users SET numLikes = numLikes + 1 WHERE username = ?", cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post.Likes++
	} else if action == "dislike" {
		if count == 0 {
			http.Error(w, "Vous ne pouvez pas unlike un post que vous n'avez pas liké", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM likes WHERE post_id = ? AND username = ?", postID, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE users SET numLikes = numLikes - 1 WHERE username = ?", cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post.Likes--
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?category=%s", post.Category), http.StatusFound)
}
