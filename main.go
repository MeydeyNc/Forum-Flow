package main

import (
	"fmt"
	DB "forum/database"
	handler "forum/handler"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const port = ":8080"

func main() {
	var db = DB.OpenDataBase()
	defer db.Close()

	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/create_post", handler.CreatePostHandler)
	http.HandleFunc("/posts", handler.PostsHandler)
	http.HandleFunc("/home", handler.HomeHandler)
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/noconnected", handler.HomeHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/update_likes", handler.MajLikesHandler)
	http.HandleFunc("/profil", handler.ProfileHandler)
	http.HandleFunc("/user/", handler.UserHandler)

	fs := http.FileServer(http.Dir("./templates/"))
	cs := http.FileServer(http.Dir("./static"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))
	http.Handle("/static/", http.StripPrefix("/static/", cs))

	fmt.Println("(http://localhost:8080/) - Serveur lancé sur le port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
