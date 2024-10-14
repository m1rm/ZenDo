package server

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "zendo/models"
)

// MockDB is a mock database implementing the necessary methods for testing
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

// TestHandleGetTodos tests the handleGetTodos function
func TestHandleGetTodos(t *testing.T) {
    mockDB := &MockDB{}
    server := &Server{Port: 8080, db: mockDB}

    req, err := http.NewRequest("GET", "/todos", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.handleGetTodos)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedTodos := []models.Todo{
        {Id: 1, Description: "Test Todo 1", Status: 0},
        {Id: 2, Description: "Test Todo 2", Status: 1},
    }
    expectedResponse, _ := json.Marshal(expectedTodos)

    if rr.Body.String() != string(expectedResponse) {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedResponse))
    }
}

// TestHandleGetTodo tests the handleGetTodo function
func TestHandleGetTodo(t *testing.T) {
    mockDB := &MockDB{}
    server := &Server{Port: 8080, db: mockDB}

    req, err := http.NewRequest("GET", "/todos/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.handleGetTodo)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedTodo := models.Todo{Id: 1, Description: "Test Todo", Status: 0}
    expectedResponse, _ := json.Marshal(expectedTodo)

    if rr.Body.String() != string(expectedResponse) {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedResponse))
    }
}

func TestHandleAddTodo(t *testing.T) {
    mockDB := &MockDB{}
    server := &Server{Port: 8080, db: mockDB}

    newTodo := models.Todo{Description: "New Test Todo", Status: 0}
    todoJSON, _ := json.Marshal(newTodo)
    req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(todoJSON))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.handleAddTodo)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
    }

    expectedResponse := `{"id":0,"description":"New Test Todo","status":0}`
    if rr.Body.String() != expectedResponse {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
    }
}

// TestHandleDeleteTodo tests the handleDeleteTodo function
func TestHandleDeleteTodo(t *testing.T) {
    mockDB := &MockDB{}
    server := &Server{Port: 8080, db: mockDB}

    req, err := http.NewRequest("DELETE", "/todos/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.handleDeleteTodo)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedResponse := `{"message":"Todo successfully deleted"}`
    if rr.Body.String() != expectedResponse {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
    }
}

// TestHandleUpdateTodo tests the handleUpdateTodo function
func TestHandleUpdateTodo(t *testing.T) {
    mockDB := &MockDB{}
    server := &Server{Port: 8080, db: mockDB}

    updatedTodo := models.Todo{Id: 1, Description: "Updated Test Todo", Status: 1}
    todoJSON, _ := json.Marshal(updatedTodo)
    req, err := http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(todoJSON))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.handleUpdateTodo)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedResponse := `{"message":"Todo successfully updated"}`
    if rr.Body.String() != expectedResponse {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
    }
}
