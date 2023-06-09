package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	// DB接続
	var err error
	db, err = sqlx.Connect("mysql", "root:secret@tcp(127.0.0.1:3306)/go_todo?parseTime=true")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	fmt.Println("server start")
	http.HandleFunc("/todos", todosHandler)

	// サーバー起動
	http.ListenAndServe(":8080", nil)
}

type Todo struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ResponseWriterがレスポンスの中身を書き込む、http.Requestでリクエストを受け取る
func todosHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストのメソッドで分岐
	switch r.Method {
	case "GET":
		getTodo(w, r)
	case "POST":
		createTodo(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	todos := []Todo{}
	err := db.Select(&todos, "SELECT id, title, content FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// wのheaderにセット
	w.Header().Set("Content-Type", "application/json")
	// todosをJSON形式に変換してwに書き込み
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo

	// JSON形式でリクエストボディをデコード
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// データベースにユーザーを追加
	res, err := db.NamedExec("INSERT INTO todos (title, content) VALUES (:title, :content)", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 追加したユーザーのIDを取得
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ID = int(id)

	// 追加したユーザーをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
