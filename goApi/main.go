package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Todo struct {
	Id          *int64 `json:"id"`
	Description string `json:"text"`
	Status      int    `json:"status"`
}

func queryTodos() ([]Todo, error) {
	var todos []Todo

	rows, err := db.Query("SELECT * FROM todos")
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

func queryTodo(id int64) (Todo, error) {
	var todo Todo

	row := db.QueryRow("SELECT * FROM todos WHERE id = ?", id)
	if err := row.Scan(&todo.Id, &todo.Description, &todo.Status); err != nil {
		if err == sql.ErrNoRows {
			return todo, fmt.Errorf("queryTodo %d: no such todo", id)
		}
		return todo, fmt.Errorf("queryTodo %d: %v", id, err)
	}
	return todo, nil
}

func insertTodo(w http.ResponseWriter, todo Todo) (int64, error) {
	query := "INSERT INTO `todos` (`description`, `status`) VALUES (?, ?);"
	insertResult, err := db.ExecContext(context.Background(), query, &todo.Description, &todo.Status)
	if err != nil {
		http.Error(w, fmt.Sprintf("error saving todo, %v", err), http.StatusInternalServerError)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not retrieve last inserted id, %v", err), http.StatusInternalServerError)
	}
	return id, nil
}

func deleteTodo(w http.ResponseWriter, todo Todo) (int64, error) {
	query := "DELETE FROM `todos` WHERE (`id`) VALUES (?);"
	insertResult, err := db.ExecContext(context.Background(), query, &todo.id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error deletomg todo, %v", err), http.StatusInternalServerError)
	}
	return id, nil
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

	todo, err := queryTodo(int64(id))
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

// @see https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func addTodo(w http.ResponseWriter, req *http.Request) {
	var t Todo
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing transmitted data, %v", err), http.StatusInternalServerError)
		return
	}
	id, err := insertTodo(w, t)
	if err != nil {
		http.Error(w, fmt.Sprintf("error saving todo in DB, %v", err), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func deleteTodo(w http.ResponseWriter, req *http.Request) {
	var t Todo
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing transmitted data, %v", err), http.StatusInternalServerError)
		return
	}
	id, err := deleteTodo(w, t)
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting todo in DB, %v", err), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(id)
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
	http.HandleFunc("POST /todos", addTodo)

	http.ListenAndServe(":8090", nil)
}
