package server

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "zendo/models"
)

// MockDB is a mock database implementing the necessary methods for testing
type MockDB struct{}

func (db *MockDB) GetTodos() ([]models.Todo, error) {
    return []models.Todo{
        {Id: 1, Description: "Test Todo 1", Status: 0},
        {Id: 2, Description: "Test Todo 2", Status: 1},
    }, nil
}


// TestHandleGetTodos tests the handleGetTodos function
func TestHandleGetTodos(t *testing.T) {
    mockDB := &MockDB{}
    server := &Server{Port: ":8080", db: mockDB}

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

// Continue writing tests for other handlers...
