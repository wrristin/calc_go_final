package orchestrator

import (
	"time"
)

type Expression struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

type Task struct {
	ID            string        `json:"id"`
	ExpressionID  string        `json:"expression_id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
	Status        string        `json:"status"`
	Result        float64       `json:"result,omitempty"`
}
