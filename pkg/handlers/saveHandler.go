package handlers

import (
	"encoding/json"
	"gomod/pkg/database"
	"gomod/pkg/entities"
	"net/http"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
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

	var userFromDB entities.UserDB
	userFromDB, err = database.GetUserByEmail(userFromRequest.Email)
	if err != nil {
		response.RespStr = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if userFromDB.Password != userFromRequest.Password {
		response.RespStr = "password is wrong"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, stock := range userFromRequest.CurAllocation {
		database.InsertStock(userFromDB.ID, stock.Ticker, stock.Number, 0)
	}

	for _, stock := range userFromRequest.TargetAllocation {
		database.InsertStock(userFromDB.ID, stock.Ticker, 0, stock.Percent)
	}

	response.RespStr = "user data is saved"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
