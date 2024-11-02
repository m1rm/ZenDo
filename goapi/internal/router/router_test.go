package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"zendo/models"
)

type MockDB struct{}

func (db *MockDB) Health() map[string]string {
	return map[string]string{"status": "healthy"}
}

func (db *MockDB) Close() error {
	// Mock close method
	return nil
}

func (db *MockDB) GetTodos() ([]models.Todo, error) {
	return []models.Todo{
		{Id: 1, Description: "Test Todo 1", Status: 0},
		{Id: 2, Description: "Test Todo 2", Status: 1},
	}, nil
}

func (db *MockDB) GetTodo(id int64) (models.Todo, error) {
	return models.Todo{Id: id, Description: "Test Todo", Status: 0}, nil
}

func (db *MockDB) DeleteTodo(id int64) (int64, error) {
	return id, nil
}

func (db *MockDB) InsertTodo(todo *models.Todo) error {
	return nil
}

func (db *MockDB) UpdateTodo(todo *models.Todo) error {
	return nil
}

func TestGetTodosHandler(t *testing.T) {
	mockDB := &MockDB{}

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	rr := NewRouter(mockDB)
	rr.ServeHTTP(w, req)

	res := w.Result()
	expected, _ := mockDB.GetTodos()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	var todos []models.Todo

	if err := json.Unmarshal(body, &todos); err != nil {
		t.Errorf("error unmarshalling response body: %v", err)
	}

	if !todosEqual(expected, todos) {
		t.Errorf("expected %v, got %v", expected, todos)
	}
}

func todosEqual(expected, actual []models.Todo) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := 0; i < len(actual); i++ {
		if actual[i].Id != expected[i].Id ||
		actual[i].Description != expected[i].Description ||
		actual[i].Status != expected[i].Status {
			return false
		}
	}
	return true
}
