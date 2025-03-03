package handlers

import (
	"database/sql"
	"encoding/json"
	"gomod/pkg/entities"
	"gomod/pkg/repository"
	"net/http"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var requestUser entities.User
	err := json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
	}

	var user entities.User
	err = repository.DB.QueryRow(`SELECT email, password, cur_allocation, target_allocation FROM users WHERE email = $1`, requestUser.Email).Scan(&user.Email, &user.Password, &user.CurAllocation, &user.TargetAllocation)
	if err != nil {
			if err == sql.ErrNoRows {
					http.Error(w, "Пользователь не найден", http.StatusNotFound)
					return
			}
			http.Error(w, "Ошибка поиска пользователя", http.StatusInternalServerError)
			return
	}

	if user.Password != requestUser.Password {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
	}

	response := struct {
			Email           string `json:"email"`
			CurAllocation   string `json:"cur_allocation"`
			TargetAllocation string `json:"target_allocation"`
			TaxRate int `json:"tax_rate"`
	}{
			Email:           user.Email,
			CurAllocation:   user.CurAllocation,
			TargetAllocation: user.TargetAllocation,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}