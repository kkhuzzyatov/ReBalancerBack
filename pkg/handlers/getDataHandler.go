package handlers

import (
	"database/sql"
	"encoding/json"
	"gomod/pkg/entities"
	"gomod/pkg/repository"
	"gomod/pkg/utils"
	"net/http"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var userOfRequest entities.User
	err := json.NewDecoder(r.Body).Decode(&userOfRequest)
	if err != nil {
		w.Write([]byte("Ошибка декодирования JSON: " + err.Error()))
		return
	}

	userOfDB, err := repository.SelectUser(userOfRequest.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Write([]byte("Не удалось загрузить данные. Проверьте, правильно ли введена почта."))
			return
		} else {
			w.Write([]byte("Ошибка сервера. Не удалось найти пользователя. Пожалуйста, попробуйте ещё раз."))
			return
		}
	}

  passwordMatches := utils.CompareHashAndPassword(userOfRequest.PasswordHash, userOfDB.PasswordHash)
	if !passwordMatches {
    w.Write([]byte("Не удалось загрузить данные. Проверьте, правильно ли введён пароль."))
		return
  }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userOfDB)
}
