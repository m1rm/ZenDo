package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"zendo/internal/database"
	"zendo/models"
)

type Router interface  {
	http.Handler
	database.Service
	registerRoutes() *http.ServeMux
}

type router struct {
	mux *http.ServeMux
	db database.Service
}

func (rr *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rr.mux.ServeHTTP(w, req)
}

func NewRouter(db database.Service) *router {
	rr := &router{
		db: db,
	}
	rr.mux = rr.registerRoutes()
	return rr
}

func (rr *router) GetMux() *http.ServeMux {
	return rr.mux
}

func (rr *router) registerRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", rr.HealthHandler)

	mux.HandleFunc("GET /todos", rr.GetTodosHandler)
	mux.HandleFunc("GET /todos/{id}", rr.GetTodoHandler)
	mux.HandleFunc("OPTIONS /todos/{id}", rr.OptionsRequestHandler)
	mux.HandleFunc("POST /todos", rr.AddTodoHandler)
	mux.HandleFunc("PUT /todos/{id}", rr.UpdateTodoHandler)
	mux.HandleFunc("DELETE /todos/{id}", rr.DeleteTodoHandler)

	return mux
}

func (rr router) HealthHandler(w http.ResponseWriter, r *http.Request) {

	jsonResp, err := json.Marshal(rr.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func successResponseHandler(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
}

func (rr router) OptionsRequestHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusNoContent)
}

func (rr router) GetTodosHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	todos, err := rr.db.GetTodos()
	if err != nil {
		http.Error(w, fmt.Sprintf("error querying the database, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	successResponseHandler(&w)
	w.Write(response)
}

func (rr router) GetTodoHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}

	todo, err := rr.db.GetTodo(int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("error querying the database, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	successResponseHandler(&w)
	w.Write(response)
}

// @see https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func (rr router) AddTodoHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	var todo models.Todo

	if err := json.NewDecoder(req.Body).Decode(&todo); err != nil {
		http.Error(w, fmt.Sprintf("error parsing transmitted data, %v", err), http.StatusInternalServerError)
		return
	}

	if err := rr.db.InsertTodo(&todo); err != nil {
		http.Error(w, fmt.Sprintf("error saving todo in DB, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&todo)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (rr router) UpdateTodoHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	var todo models.Todo
	if err := json.NewDecoder(req.Body).Decode(&todo); err != nil {
		http.Error(w, fmt.Sprintf("error parsing transmitted data, %v", err), http.StatusInternalServerError)
		return
	}

	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}

	todo.Id = int64(id) // Ensure the ID is set correctly
	if err := rr.db.UpdateTodo(&todo); err != nil {
		http.Error(w, fmt.Sprintf("error updating todo in DB, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&todo)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	successResponseHandler(&w)
	w.Write(response)
}

func (rr router) DeleteTodoHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}
	_, err = rr.db.DeleteTodo(int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting todo in DB, %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
