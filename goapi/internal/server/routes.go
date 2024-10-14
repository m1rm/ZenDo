package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"zendo/models"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)

	mux.HandleFunc("GET /todos", s.handleGetTodos)
	mux.HandleFunc("GET /todos/{id}", s.handleGetTodo)
	mux.HandleFunc("OPTIONS /todos/{id}", handleOptionsRequest)
	mux.HandleFunc("POST /todos", s.handleAddTodo)
	mux.HandleFunc("PUT /todos/{id}", s.handleUpdateTodo)
	mux.HandleFunc("DELETE /todos/{id}", s.handleDeleteTodo)

	return mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

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

func handleSuccessResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
}

func handleOptionsRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetTodos(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	todos, err := s.db.GetTodos()
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

func (s *Server) handleGetTodo(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}

	todo, err := s.db.GetTodo(int64(id))
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
func (s *Server) handleAddTodo(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	var todo models.Todo

	if err := json.NewDecoder(req.Body).Decode(&todo); err != nil {
		http.Error(w, fmt.Sprintf("error parsing transmitted data, %v", err), http.StatusInternalServerError)
		return
	}

	if err := s.db.InsertTodo(&todo); err != nil {
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

func (s *Server) handleUpdateTodo(w http.ResponseWriter, req *http.Request) {
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
	if err := s.db.UpdateTodo(&todo); err != nil {
		http.Error(w, fmt.Sprintf("error updating todo in DB, %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&todo)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	handleSuccessResponse(&w)
	w.Write(response)
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	queryParam := req.PathValue("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing id query param, %v", err), http.StatusInternalServerError)
		return
	}
	_, err = s.db.DeleteTodo(int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting todo in DB, %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
