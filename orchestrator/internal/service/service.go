package service

import (
	"calc_service/models"
	"calc_service/orchestrator/utils"
	"errors"
	"sync"
)

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrTaskNotFound      = errors.New("task not found")
)

var (
	expressions = make(map[string]models.Expression)
	tasks       = make(chan models.Task, 100)
	taskStore   = make(map[string]models.Task)
	mu          sync.Mutex
)

func AddExpression(expr string) (string, error) {
	id := utils.GenerateID()
	exp := models.Expr{
		ID:   id,
		Expr: expr,
	}

	newTasks, err := utils.ParseExpression(exp)
	if err != nil {
		return "", ErrInvalidExpression
	}

	mu.Lock()
	expressions[id] = models.Expression{
		ID:     id,
		Expr:   expr,
		Status: "pending",
	}

	for _, task := range newTasks {
		task.ExpressionID = id
		taskStore[task.ID] = task
		tasks <- task
	}
	mu.Unlock()

	return id, nil
}

func GetTask() (models.Task, error) {
	for {
		select {
		case task := <-tasks:
			// Проверяем, готова ли задача к выполнению
			if task.Status == "waiting" {
				if (task.Arg1TaskID != "" && !IsTaskCompleted(task.Arg1TaskID)) ||
					(task.Arg2TaskID != "" && !IsTaskCompleted(task.Arg2TaskID)) {
					// Если зависимая задача ещё не готова, возвращаем её обратно в очередь
					tasks <- task
					continue
				}
			}

			// Если зависимые задачи готовы, изменяем статус и отправляем задачу агенту
			task.Status = "in_progress"
			return task, nil

		default:
			return models.Task{}, ErrTaskNotFound
		}
	}
}

func GetTaskByID(taskID string) (models.Task, error) {
	task, exists := taskStore[taskID] // Ищем задачу в хранилище
	if !exists {
		return models.Task{}, errors.New("task not found") // Если нет – возвращаем ошибку
	}
	return task, nil // Возвращаем найденную задачу
}

func IsTaskCompleted(taskID string) bool {
	task, err := GetTaskByID(taskID) // Получаем задачу по ID
	if err != nil {
		return false
	}
	return task.Status == "done"
}

func CompleteTask(taskID string, result float64) error {
	mu.Lock()
	defer mu.Unlock()

	task, exists := taskStore[taskID]
	if !exists {
		return ErrTaskNotFound
	}

	task.Result = result
	task.Status = "completed"
	taskStore[taskID] = task

	expr, exists := expressions[task.ExpressionID]
	if !exists {
		return ErrTaskNotFound
	}

	allCompleted := true
	var finalResult float64

	for _, t := range taskStore {
		if t.ExpressionID == expr.ID {
			if t.Status != "completed" {
				allCompleted = false
				break
			}
			finalResult = t.Result
		}
	}

	if allCompleted {
		expr.Result = finalResult
		expr.Status = "done"
		expressions[task.ExpressionID] = expr
	}

	return nil
}

func GetExpressions() []models.Expression {
	mu.Lock()
	defer mu.Unlock()

	exprList := make([]models.Expression, 0, len(expressions))
	for _, expr := range expressions {
		exprList = append(exprList, models.Expression{
			ID:     expr.ID,
			Status: expr.Status,
			Result: expr.Result,
		})
	}
	return exprList
}

func GetExpressionByID(id string) (models.Expression, error) {
	mu.Lock()
	defer mu.Unlock()

	expr, exists := expressions[id]
	if !exists {
		return models.Expression{}, errors.New("expression not found")
	}

	return expr, nil
}
