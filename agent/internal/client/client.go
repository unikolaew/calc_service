package client

import (
	"bytes"
	"calc_service/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func GetTask() (models.Task, error) {
	resp, err := http.Get("http://localhost/internal/task")
	if err != nil {
		return models.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.Task{}, errors.New("no tasks available")
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func PostResult(taskID string, result float64) error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"id":     taskID,
		"result": result,
	})

	resp, err := http.Post("http://localhost/internal/task", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%d", resp.StatusCode))
	}

	return nil
}
