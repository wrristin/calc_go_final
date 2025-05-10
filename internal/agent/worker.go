package agent

import (
	"bytes"
	"calc_service/internal/orchestrator"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

func Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task := getTaskGRPC() // Используем gRPC вместо HTTP
		if task != nil {
			result := executeTask(task)
			sendResultGRPC(task.ID, result) // Используем gRPC
		}
		time.Sleep(1 * time.Second)
	}
}

func getTask() *orchestrator.Task {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		log.Println("Failed to get task:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("No tasks available, status code: %d", resp.StatusCode)
		return nil
	}

	var task orchestrator.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		log.Println("Failed to decode task:", err)
		return nil
	}
	return &task
}

func executeTask(task *orchestrator.Task) float64 {
	log.Printf("Executing task: ID=%s, Arg1=%f, Arg2=%f, Operation=%s", task.ID, task.Arg1, task.Arg2, task.Operation)
	time.Sleep(task.OperationTime)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}

func sendResult(taskID string, result float64) {
	log.Printf("Sending result for task: ID=%s, Result=%f", taskID, result)
	payload := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to send result:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send result, status code: %d", resp.StatusCode)
		return
	}
	log.Println("Result sent for task ID:", taskID)
}
