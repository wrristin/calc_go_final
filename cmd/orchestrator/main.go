package main

import (
	"calc_service/internal/orchestrator"
	"log"
	"net/http"
)

func main() {
	// Инициализация БД
	if err := orchestrator.InitDB(); err != nil {
		log.Fatal("DB init failed:", err)
	}

	// Запуск gRPC сервера
	go orchestrator.StartGRPCServer()

	// HTTP роуты
	http.HandleFunc("/api/v1/register", orchestrator.RegisterHandler)
	http.HandleFunc("/api/v1/login", orchestrator.LoginHandler)
	http.HandleFunc("/api/v1/calculate", orchestrator.AuthMiddleware(orchestrator.AddExpressionHandler))
	http.HandleFunc("/api/v1/expressions", orchestrator.AuthMiddleware(orchestrator.GetExpressionsHandler))
	http.HandleFunc("/api/v1/expressions/", orchestrator.AuthMiddleware(orchestrator.GetExpressionByIDHandler))

	log.Println("Starting orchestrator on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
