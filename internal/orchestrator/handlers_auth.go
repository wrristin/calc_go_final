package orchestrator

import (
	"calc_service/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword := auth.HashPassword(req.Password)
	user := User{
		ID:       uuid.New().String(),
		Login:    req.Login,
		Password: hashedPassword,
	}

	if err := DB.CreateUser(user); err != nil {
		http.Error(w, "User exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := DB.GetUserByLogin(req.Login)
	if err != nil || auth.HashPassword(req.Password) != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
