package handler

import (
	"calc_service/orchestrator/internal/service"
	"encoding/json"
	"net/http"
)

// Добавление выражения
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := service.AddExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// Получение списка выражений
func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions := service.GetExpressions()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

// Получение выражения по ID
func HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	expr, err := service.GetExpressionByID(id)
	if err != nil {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

// Обработка задач (для агентов)
func HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		task, err := service.GetTask()
		if err != nil {
			http.Error(w, "No tasks available", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	case http.MethodPost:
		var result struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := service.CompleteTask(result.ID, result.Result); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
