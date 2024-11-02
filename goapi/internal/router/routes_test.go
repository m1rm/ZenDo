package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"zenod/models"
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
	rr := router.NewRouter(mockDB)

	w := httptest.NewRecorder()
 	handler := http.HandlerFunc(rr.GetTodosHandler)
  	handler.ServeHTTP(w, req)

	res := w.Result()


}
