package repository

import (
	"database/sql"
	"fmt"
	"gomod/config"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("postgres", config.GetDBConnStr())
    if err != nil {
        log.Fatalf("Ошибка подключения к базе данных: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Ошибка проверки подключения: %v", err)
    }

    createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				password VARCHAR(255) NOT NULL,
				cur_allocation TEXT,
				target_allocation TEXT
    );`
    _, err = DB.Exec(createTableQuery)
    if err != nil {
        log.Fatalf("Ошибка создания таблицы: %v", err)
    }

    fmt.Println("База данных успешно подключена и таблица создана.")
}