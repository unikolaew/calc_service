package handler_test

import (
	"bytes"
	"calc_service/orchestrator/internal/handler"
	"calc_service/orchestrator/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCalculate(t *testing.T) {
	reqBody := map[string]string{"expression": "5 + 5"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCalculate(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestHandleGetExpressions(t *testing.T) {
	service.AddExpression("10 / 2")

	req := httptest.NewRequest("GET", "/api/v1/expressions", nil)
	w := httptest.NewRecorder()

	handler.HandleGetExpressions(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
