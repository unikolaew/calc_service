package service_test

import (
	"calc_service/orchestrator/internal/service"
	"testing"
)

func TestAddExpression(t *testing.T) {
	expr := "2 + 2"
	id, err := service.AddExpression(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestGetTask(t *testing.T) {
	_, _ = service.AddExpression("3 * 3")
	task, err := service.GetTask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.ID == "" {
		t.Fatal("expected non-empty task ID")
	}
}
