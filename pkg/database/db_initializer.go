package database

import (
	"database/sql"
	"fmt"
	"gomod/config"
)

var DB *sql.DB

func InitDB() error {
	var err error

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetDatabaseHost(), config.GetDatabasePort(), config.GetDatabaseUser(), config.GetDatabasePassword(), config.GetDatabaseName())

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка при проверке подключения: %v", err)
	}

	return nil
}
