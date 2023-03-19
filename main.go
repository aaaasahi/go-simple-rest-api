package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	// DB接続
	var err error
	db, err = sqlx.Open("mysql", "root:secret@tcp(127.0.0.1:3306)/go_todo?parseTime=true")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Database connection successful!")
}
