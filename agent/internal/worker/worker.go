package worker

import (
	"calc_service/agent/internal/client"
	"calc_service/models"
	"log"
	"os"
	"strconv"
	"time"
)

func Start() {
	for {
		task, err := client.GetTask()
		if err != nil {
			log.Println("Failed to get task:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		zero := map[string]models.Task{}
		result := computeTask(task, zero)
		if err := client.PostResult(task.ID, result); err != nil {
			log.Println("Failed to post result:", err)
		} else {
			log.Printf("Task %s completed with result: %f\n", task.ID, result)
		}
	}
}

func computeTask(task models.Task, taskMap map[string]models.Task) float64 {
	log.Printf("Computing task: ID=%s, Operation=%s, Arg1=%f, Arg2=%f\n", task.ID, task.Operation, task.Arg1, task.Arg2)

	// Загружаем время выполнения операций из переменных окружения
	additionTime := getDurationFromEnv("TIME_ADDITION_MS", 100)
	subtractionTime := getDurationFromEnv("TIME_SUBTRACTION_MS", 100)
	multiplicationTime := getDurationFromEnv("TIME_MULTIPLICATIONS_MS", 100)
	divisionTime := getDurationFromEnv("TIME_DIVISIONS_MS", 100)

	// Получаем реальные аргументы (если они идут из предыдущих задач)
	arg1 := task.Arg1
	arg2 := task.Arg2

	if task.Arg1TaskID != "" {
		prevTask, exists := taskMap[task.Arg1TaskID]
		if exists {
			arg1 = computeTask(prevTask, taskMap)
		} else {
			log.Printf("Task ID not found: %s", task.Arg1TaskID)
		}
	}

	if task.Arg2TaskID != "" {
		prevTask, exists := taskMap[task.Arg2TaskID]
		if exists {
			arg2 = computeTask(prevTask, taskMap)
		} else {
			log.Printf("Task ID not found: %s", task.Arg2TaskID)
		}
	}

	// Вычисляем результат
	var result float64
	switch task.Operation {
	case "+":
		time.Sleep(additionTime)
		result = arg1 + arg2
	case "-":
		time.Sleep(subtractionTime)
		result = arg1 - arg2
	case "*":
		time.Sleep(multiplicationTime)
		result = arg1 * arg2
	case "/":
		if arg2 == 0 {
			log.Println("Division by zero in task", task.ID)
			return 0
		}
		time.Sleep(divisionTime)
		result = arg1 / arg2
	default:
		log.Println("Unknown operation:", task.Operation)
		return 0
	}

	log.Printf("Computed task: ID=%s, Result=%f\n", task.ID, result)
	return result
}

// getDurationFromEnv загружает значение из переменной окружения и возвращает его как time.Duration.
// Если переменная окружения не установлена, возвращается значение по умолчанию.
func getDurationFromEnv(envVar string, defaultValue int) time.Duration {
	valStr := os.Getenv(envVar)
	if valStr == "" {
		return time.Duration(defaultValue) * time.Millisecond
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid value for %s: %v. Using default value: %d ms\n", envVar, err, defaultValue)
		return time.Duration(defaultValue) * time.Millisecond
	}

	return time.Duration(val) * time.Millisecond
}
