package orchestrator

import (
	"calc_service/pkg/database"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

var DB *Repository

func InitDB() error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	DB = &Repository{db: db}
	return nil
}

// Методы для работы с пользователями
func (r *Repository) CreateUser(user User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (id, login, password) VALUES (?, ?, ?)",
		user.ID, user.Login, user.Password,
	)
	return err
}

func (r *Repository) GetUserByLogin(login string) (*User, error) {
	row := r.db.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)
	user := &User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	return user, err
}

// Методы для работы с выражениями
func (r *Repository) CreateExpression(expr Expression) error {
	_, err := r.db.Exec(
		"INSERT INTO expressions (id, user_id, status, result) VALUES (?, ?, ?, ?)",
		expr.ID, expr.UserID, expr.Status, expr.Result,
	)
	return err
}

func (r *Repository) GetExpressionsByUser(userID string) ([]Expression, error) {
	rows, err := r.db.Query("SELECT id, status, result FROM expressions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []Expression
	for rows.Next() {
		var expr Expression
		if err := rows.Scan(&expr.ID, &expr.Status, &expr.Result); err != nil {
			return nil, err
		}
		expr.UserID = userID // Восстанавливаем UserID из контекста
		expressions = append(expressions, expr)
	}
	return expressions, nil
}

func (r *Repository) GetExpressionByID(id string) (*Expression, error) {
	row := r.db.QueryRow("SELECT id, user_id, status, result FROM expressions WHERE id = ?", id)
	expr := &Expression{}
	err := row.Scan(&expr.ID, &expr.UserID, &expr.Status, &expr.Result)
	return expr, err
}

// Методы для работы с задачами
func (r *Repository) CreateTask(task Task) error {
	_, err := r.db.Exec(
		`INSERT INTO tasks 
		(id, expression_id, arg1, arg2, operation, operation_time, status, result) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.ExpressionID, task.Arg1, task.Arg2,
		task.Operation, task.OperationTime.Milliseconds(), task.Status, task.Result,
	)
	return err
}

func (r *Repository) GetPendingTask() (*Task, error) {
	row := r.db.QueryRow(`
		SELECT id, expression_id, arg1, arg2, operation, operation_time 
		FROM tasks 
		WHERE status = 'pending' 
		LIMIT 1
	`)
	task := &Task{}
	var opTime int64
	err := row.Scan(
		&task.ID, &task.ExpressionID, &task.Arg1, &task.Arg2,
		&task.Operation, &opTime,
	)
	task.OperationTime = time.Duration(opTime) * time.Millisecond
	task.Status = "pending" // Статус не выбирается из БД, но явно задаем
	return task, err
}

func (r *Repository) UpdateTaskResult(taskID string, result float64) error {
	_, err := r.db.Exec(
		"UPDATE tasks SET result = ?, status = 'completed' WHERE id = ?",
		result, taskID,
	)
	return err
}

func (r *Repository) UpdateTaskStatus(taskID string, status string) error {
	_, err := r.db.Exec(
		"UPDATE tasks SET status = ? WHERE id = ?",
		status, taskID,
	)
	return err
}

func (r *Repository) GetExpressionByTaskID(taskID string) (*Expression, error) {
	row := r.db.QueryRow(`
		SELECT e.id, e.user_id, e.status, e.result 
		FROM expressions e
		JOIN tasks t ON e.id = t.expression_id
		WHERE t.id = ?`, taskID)

	expr := &Expression{}
	err := row.Scan(&expr.ID, &expr.UserID, &expr.Status, &expr.Result)
	return expr, err
}

func (r *Repository) CheckAllTasksCompleted(expressionID string) (bool, error) {
	rows, err := r.db.Query(
		"SELECT status FROM tasks WHERE expression_id = ?",
		expressionID,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			return false, err
		}
		if status != "completed" {
			return false, nil
		}
	}
	return true, nil
}

func (r *Repository) UpdateExpressionStatus(expressionID, status, result string) error {
	_, err := r.db.Exec(
		"UPDATE expressions SET status = ?, result = ? WHERE id = ?",
		status, result, expressionID,
	)
	return err
}
