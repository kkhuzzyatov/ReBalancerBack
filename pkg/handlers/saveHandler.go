package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gomod/pkg/entities"
	"gomod/pkg/repository"
	"gomod/pkg/stocks"
	"gomod/pkg/utils"
	"net/http"
)

func Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
	}

	err = stocks.IsAllocValid[int](utils.AllocationParser[int](user.CurAllocation)) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = stocks.IsAllocValid[float64](utils.AllocationParser[float64](user.TargetAllocation)) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var storedPassword string
	err = repository.DB.QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&storedPassword)
	if err != nil {
			if err == sql.ErrNoRows {

				err = repository.DB.QueryRow(`INSERT INTO users (email, password, cur_allocation, target_allocation) VALUES ($1, $2, $3, $4) RETURNING email`, user.Email, user.Password, user.CurAllocation, user.TargetAllocation).Scan(&user.Email)
				if err != nil {
					http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
					return
				}
	
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("Пользователь с email %s успешно создан", user.Email)))
				return
			}
			http.Error(w, "Ошибка проверки пользователя", http.StatusInternalServerError)
			return
	}

	if storedPassword != user.Password {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
	}

	_, err = repository.DB.Exec(`UPDATE users SET cur_allocation = $3, target_allocation = $4 WHERE email = $1 AND password = $2`, user.Email, user.Password, user.CurAllocation, user.TargetAllocation)
	if err != nil {
			http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}