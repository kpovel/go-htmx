package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type App struct {
	DB *sql.DB
}

type Todo struct {
	id        int
	title     string
	completed bool
}

func dbClient() *sql.DB {
	var dbUrl = "file:///tmp/go-htmx-todo.db"
	db, err := sql.Open("libsql", dbUrl)

	fmt.Printf("db: %s\n", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}

	db.Exec("create table if not exists todos (id integer auto_increment primary key, title varchar(255), completed bool)")

	return db
}

func main() {
	client := dbClient()

	app := &App{
		DB: client,
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/todo", app.CreateTodo)
	http.HandleFunc("/todos", app.GetTodos)

	port := "6969"
	fmt.Printf("Listening at localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (app *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	formValue := r.FormValue("title")

	fmt.Printf("formValue: %s\n", formValue)
	app.DB.Exec("insert into todos (title) values (?)", formValue)

	fmt.Fprintf(w, `<li>%s</li>`, formValue)
}

func (app *App) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.DB.Query("select * from todos order by id desc")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get todos: %s", err)
		os.Exit(1)
	}
	defer todos.Close()

	formattedTodos := ""
	for todos.Next() {
		var id int
		var title string
		var completed bool

		err = todos.Scan(&id, &title, &completed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan db row: %s", err)
			os.Exit(1)
		}

		formattedTodos += fmt.Sprintf("<li id=\"%d\">%s</li>", id, title)
	}

  fmt.Fprintf(w, formattedTodos)
}
