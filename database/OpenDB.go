package DB

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDataBase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Println(err)
	}
	return db
}
