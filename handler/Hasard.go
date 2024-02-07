package forum

import (
	"database/sql"
)

type User struct {
	ID          sql.NullInt64
	Username    string
	Email       string
	Password    string
	NumPosts    int
	NumLikes    int
	NumDislikes int
}

type Post struct {
	ID       int
	Title    string
	Content  string
	Category string
	Likes    int
}
