package handlers

import (
	"database/sql"
	"encoding/json"
	"gomod/pkg/entities"
	"gomod/pkg/repository"
	"gomod/pkg/utils"
	"net/http"
)

func Save(w http.ResponseWriter, r *http.Request) {
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
		// Пользователь не найден по email -> создаём новый аккаунт
		if err == sql.ErrNoRows {
			userOfRequest.PasswordHash, err = utils.HashPassword(userOfRequest.PasswordHash)
			if err != nil {
				w.Write([]byte("Ошибка сервера. Что-то пошло не так во время хеширования пароля. Пожалуйста, попробуйте ещё раз."))
				return
			}

			err := repository.InsertUser(userOfRequest)
			if err != nil {
				w.Write([]byte("Ошибка сервера. Не удалось создать пользователя. Пожалуйста, попробуйте ещё раз."))
				return
			}

			w.Write([]byte("Вы успешно зарегистрированы."))
			return
		} else {
			w.Write([]byte("Ошибка сервера. Не удалось найти пользователя. Пожалуйста, попробуйте ещё раз."))
			return
		}
	}
	// Пользователь найден по email -> проверка пароля
	passwordMatches := utils.CompareHashAndPassword(userOfRequest.PasswordHash, userOfDB.PasswordHash)
	if !passwordMatches {
		w.Write([]byte("Не удалось сохранить данные. Проверьте, правильно ли введён пароль."))
		return
	}

	err = repository.UpdateUser(userOfRequest)
	if err != nil {
		w.Write([]byte("Ошибка сервера. Не удалось обновить данные. Пожалуйста, попробуйте ещё раз."))
		return
	}

	w.Write([]byte("Данные успешно обновлены."))
	return
}
