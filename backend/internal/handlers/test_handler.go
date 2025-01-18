// handlers/test_handler.go
package handlers

import (
    "context"
    "encoding/json"
    "net/http"

    "backend/internal/repository"
)

type TestHandler struct {
    testRepo *repository.TestRepository
}

func NewTestHandler(testRepo *repository.TestRepository) *TestHandler {
    return &TestHandler{testRepo: testRepo}
}

func (h *TestHandler) InsertTestData(w http.ResponseWriter, r *http.Request) {
    var payload interface{}
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }
    if err := h.testRepo.InsertTestData(context.Background(), payload); err != nil {
        http.Error(w, "Failed to insert data", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

func (h *TestHandler) GetLastTestData(w http.ResponseWriter, r *http.Request) {
    data, err := h.testRepo.GetLastTestData(context.Background())
    if err != nil {
        http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(data)
}