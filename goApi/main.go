package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	id     int32
	text   string
	status string
}

func todos(w http.ResponseWriter, req *http.Request) {
	todos := []Todo{
		{1, "Style mock content", "done"},
		{2, "Style nav", "done"},
		{3, "Add header and footer", "done"},
		{4, "Add CTAs without function", "done"},
		{5, "Add Go Api with mock responses", "open"},
		{6, "dockerize", "open"},
		{7, "add just & docker compose setup", "open"},
		{8, "add DB", "open"},
		{9, "replace mock API responses data with real data", "open"},
		{10, "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea", "open"},
	}

	response, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("GET /todos", todos)

	http.ListenAndServe(":8090", nil)
}
