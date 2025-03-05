package main

import (
	"calc_service/orchestrator/internal/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handler.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", handler.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/", handler.HandleGetExpressionByID)
	http.HandleFunc("/internal/task", handler.HandleTask)

	log.Println("Orchestrator started on localhost")
	log.Fatal(http.ListenAndServe("", nil))
}
