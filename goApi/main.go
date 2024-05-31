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
	Id          int     `json:"id"`
	Description string  `json:"text"`
	Status      int     `json:"status"`
}

func queryTodos() ([]Todo, error) {
    var todos []Todo

    rows, err := db.Query("SELECT * FROM todo")
    if err != nil {
        return nil, fmt.Errorf("queryTodos: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.Id, &todo.Description, &todo.Status); err != nil {
            return nil, fmt.Errorf("queryTodos: %v", err)
        }
        todos = append(todos, todo)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("queryTodos: %v", err)
    }
    return todos, nil
}


func queryTodo(id int) (Todo, error) {
    var todo Todo

    row := db.QueryRow("SELECT * FROM todo WHERE id = ?", id)
    if err := row.Scan(&todo.Id, &todo.Description, &todo.Status); err != nil {
        if err == sql.ErrNoRows {
            return todo, fmt.Errorf("queryTodo %d: no such todo", id)
        }
        return todo, fmt.Errorf("queryTodo %d: %v", id, err)
    }
    return todo, nil
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

    todo, err := queryTodo(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("error querying the database, %v", err), http.StatusInternalServerError)
        return
    }

	response, err := json.Marshal(todo)
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
