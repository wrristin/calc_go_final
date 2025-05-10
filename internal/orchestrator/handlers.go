package orchestrator

import (
	"calc_service/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add expression")

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Bad request:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Println("Expression received:", req.Expression)

	taskList := parseExpression(req.Expression)
	if taskList == nil {
		log.Println("Failed to parse expression:", req.Expression)
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	// Сохраняем выражение и задачи в БД
	userLogin := r.Context().Value("userLogin").(string)
	user, err := DB.GetUserByLogin(userLogin)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	expressionID := taskList[0].ExpressionID
	expression := Expression{
		ID:     expressionID,
		UserID: user.ID,
		Status: "pending",
	}
	if err := DB.CreateExpression(expression); err != nil {
		log.Printf("Failed to save expression: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	for _, task := range taskList {
		if err := DB.CreateTask(task); err != nil {
			log.Printf("Failed to save task: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": expressionID})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received request to get a task")
		task, err := DB.GetPendingTask()
		if err != nil || task == nil {
			log.Println("No tasks available")
			http.Error(w, "No tasks available", http.StatusNotFound)
			return
		}

		// Обновляем статус задачи на "in progress"
		if err := DB.UpdateTaskStatus(task.ID, "in progress"); err != nil {
			log.Printf("Failed to update task status: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		log.Printf("Task assigned: ID=%s, ExpressionID=%s", task.ID, task.ExpressionID)
		json.NewEncoder(w).Encode(task)

	case http.MethodPost:
		log.Println("Received request to submit task result")
		var result struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			log.Println("Invalid request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Обновляем результат задачи
		if err := DB.UpdateTaskResult(result.ID, result.Result); err != nil {
			log.Printf("Failed to update task result: %v", err)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// Проверяем статус выражения
		expression, err := DB.GetExpressionByTaskID(result.ID)
		if err != nil {
			log.Printf("Expression not found for task ID=%s: %v", result.ID, err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		allCompleted, err := DB.CheckAllTasksCompleted(expression.ID)
		if err != nil {
			log.Printf("Failed to check tasks: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if allCompleted {
			if err := DB.UpdateExpressionStatus(expression.ID, "completed", fmt.Sprintf("%.2f", result.Result)); err != nil {
				log.Printf("Failed to update expression: %v", err)
			}
			log.Printf("Expression completed: ID=%s, Result=%s", expression.ID, expression.Result)
		}

		w.WriteHeader(http.StatusOK)
	default:
		log.Println("Method not allowed:", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Listing all expressions")
	userLogin := r.Context().Value("userLogin").(string)
	user, err := DB.GetUserByLogin(userLogin)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	expressions, err := DB.GetExpressionsByUser(user.ID)
	if err != nil {
		log.Printf("Failed to fetch expressions: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	log.Println("Received request to get expression by ID:", id)

	expression, err := DB.GetExpressionByID(id)
	if err != nil {
		log.Println("Expression not found:", id)
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": expression,
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userLogin", claims.Login)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
