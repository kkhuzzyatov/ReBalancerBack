package handlers

import (
	"encoding/json"
	"gomod/pkg/database"
	"gomod/pkg/entities"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var response entities.UserResponse
	var userFromRequest entities.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userFromRequest)
	if err != nil {
		response.RespStr = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	err = database.InsertUser(userFromRequest.Email, userFromRequest.Password)
	if err != nil {
		response.RespStr = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response.RespStr = "user is created"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
