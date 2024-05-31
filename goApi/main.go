package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
    "database/sql"
    "log"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Todo struct {
	Id          int32  `json:"id"`
	Description string `json:"text"`
	Status      string `json:"status"`
}

func queryTodos() ([]Todo, error) {
    var todos []Todo

    rows, err := db.Query("SELECT * FROM todo")
    if err != nil {
        return nil, fmt.Errorf("todos: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.Id, &todo.Description, &todo.Status); err != nil {
            return nil, fmt.Errorf("todos: %v", err)
        }
        todos = append(todos, todo)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("todos: %v", err)
    }
    return todos, nil
}

func getTodos(w http.ResponseWriter, req *http.Request) {

    todos, err := queryTodos()
    if err != nil {
        http.Error(w, fmt.Sprintf("error querying the database, %v", err), http.StatusInternalServerError)
        return
    }

	response, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func getTodo(w http.ResponseWriter, req *http.Request) {
	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}

	// bc. array index starts at 0
	correctedId := id - 1

	todos := []Todo{
		{1, "Style mock content", "done"},
        		{2, "Style nav", "done"},
        		{3, "Add header and footer", "done"},
        		{4, "Add CTAs without function", "done"},
        		{5, "Add Go Api with mock responses", "done"},
        		{6, "dockerize dev", "done"},
        		{7, "add just & docker compose setup", "done"},
        		{8, "add DB", "done"},
        		{9, "replace mock API responses data with real data", "open"},
        		{10, "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea", "open"},
	}

	response, err := json.Marshal(todos[correctedId])
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
    // --- START Connect to DB ---
    // Capture connection properties.
    cfg := mysql.Config{
        User:   "root",
        Passwd: "ChangeMe",
        Net:    "tcp",
        Addr:   "db:3306",
        DBName: "todoApp",
    }
    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("DB Connected!")
    // --- END Connect to DB ---

	http.HandleFunc("GET /todos", getTodos)
	http.HandleFunc("GET /todos/{id}", getTodo)

	http.ListenAndServe(":8090", nil)
}
