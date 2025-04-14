package database

import (
	"fmt"
	"gomod/pkg/entities"
)

func CreateUsersTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы users: %v", err)
	}
	return nil
}

func InsertUser(email string, password string) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	_, err := DB.Exec(query, email, password)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении пользователя: %v", err)
	}
	return nil
}

func GetUserByEmail(email string) (entities.UserDB, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`
	var user entities.UserDB
	err := DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return user, fmt.Errorf("ошибка при получении пользователя: %v", err)
	}
	return user, nil
}

func DeleteUser(email string) error {
	query := `DELETE FROM users WHERE email = $1`
	_, err := DB.Exec(query, email)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %v", err)
	}
	return nil
}
