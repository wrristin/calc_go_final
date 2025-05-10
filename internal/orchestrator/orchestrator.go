package orchestrator

import (
	"calc_service/internal/calculation"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func generateID() string {
	return uuid.New().String()
}

func parseExpression(expr string) []Task {
	log.Println("Parsing expression:", expr)
	tokens := calculation.Tokenize(expr)
	postfix, err := calculation.InfixToPostfix(tokens)
	if err != nil {
		log.Println("Error parsing expression:", err)
		return nil
	}

	expressionID := generateID()
	log.Printf("Generated ExpressionID: %s", expressionID)

	var taskList []Task
	var stack []float64

	for _, token := range postfix {
		if calculation.IsNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if calculation.IsOperator(token) {
			if len(stack) < 2 {
				log.Println("Not enough operands for operator:", token)
				return nil
			}
			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			task := Task{
				ID:            generateID(),
				Arg1:          arg1,
				Arg2:          arg2,
				Operation:     token,
				OperationTime: getOperationTime(token),
				Status:        "pending",
				ExpressionID:  expressionID,
			}
			taskList = append(taskList, task)
			log.Printf("Created task: ID=%s, ExpressionID=%s", task.ID, task.ExpressionID)
		}
	}

	// Сохраняем выражение в БД
	expression := Expression{
		ID:     expressionID,
		Status: "pending",
	}
	if err := DB.CreateExpression(expression); err != nil {
		log.Printf("Failed to save expression: %v", err)
		return nil
	}

	// Сохраняем задачи в БД
	for _, task := range taskList {
		if err := DB.CreateTask(task); err != nil {
			log.Printf("Failed to save task: %v", err)
			return nil
		}
	}

	log.Printf("Parsed expression: ID=%s, Tasks=%d", expressionID, len(taskList))
	return taskList
}

func getOperationTime(op string) time.Duration {
	var timeMs int
	var err error

	switch op {
	case "+":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
		if err != nil {
			timeMs = 1000
		}
	case "-":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
		if err != nil {
			timeMs = 1000
		}
	case "*":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
		if err != nil {
			timeMs = 2000
		}
	case "/":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
		if err != nil {
			timeMs = 2000
		}
	default:
		timeMs = 1000
	}

	log.Printf("Operation %s time: %d ms", op, timeMs)
	return time.Duration(timeMs) * time.Millisecond
}
