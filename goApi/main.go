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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleSuccessResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
}

func handleOptionsRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusNoContent)
}

func getTodos() ([]Todo, error) {
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

func getTodo(id int64) (Todo, error) {
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

func deleteTodo(w http.ResponseWriter, id int64) (int64, error) {
	_, err := getTodo(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error retrieving todo from database, %v", err), http.StatusInternalServerError)
		return 0, err
	}

	query := "DELETE FROM `todos` WHERE `id` = ?;"

	_, execErr := db.ExecContext(context.Background(), query, id)
	if execErr != nil {
		http.Error(w, fmt.Sprintf("error deleting todo, %v", err), http.StatusInternalServerError)
		return 0, err
	}

	return id, nil
}

func handleGetTodos(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	todos, err := getTodos()
	if err != nil {
		http.Error(w, fmt.Sprintf("error querying the database, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(&w)
	w.Write(response)
}

func handleGetTodo(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}

	todo, err := getTodo(int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("error querying the database, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(&w)
	w.Write(response)
}

// @see https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func handleAddTodo(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

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

	handleSuccessResponse(&w)
	w.Write(response)
}

func handleDeleteTodo(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	var id int64
	err := json.NewDecoder(req.Body).Decode(&id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing transmitted data, %v", err), http.StatusInternalServerError)
		return
	}
	id, err = deleteTodo(w, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting todo in DB, %v", err), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(&w)
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

	http.HandleFunc("GET /todos", handleGetTodos)
	http.HandleFunc("GET /todos/{id}", handleGetTodo)
	http.HandleFunc("OPTIONS /todos/{id}", handleOptionsRequest)
	http.HandleFunc("POST /todos", handleAddTodo)
	http.HandleFunc("DELETE /todos/{id}", handleDeleteTodo)

	http.ListenAndServe(":8090", nil)
}
