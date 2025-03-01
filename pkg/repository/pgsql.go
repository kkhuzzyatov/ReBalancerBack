package repository

import (
	"database/sql"
	"gomod/config"
	"gomod/pkg/entities"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func SetDB(db *sql.DB) {
	db = db
}

func InitDB() {
	var err error

	db, err = sql.Open("postgres", config.GetDBConnStr())
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to reach the database: %v", err)
	}
}

func SelectUser(email string) (entities.User, error) {
	var user entities.User
	err := db.QueryRow("SELECT id, email, password, curAllocation, targetAllocation, taxRate FROM users WHERE email=$1", email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CurAllocation, &user.TargetAllocation, &user.TaxRate)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func InsertUser(user entities.User) error {
	var id int
	err := db.QueryRow("INSERT INTO users (email, password, curAllocation, targetAllocation, taxRate) VALUES ($1, $2, $3, $4, $5) RETURNING id", user.Email, user.PasswordHash, user.CurAllocation, user.TargetAllocation, user.TaxRate).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(user entities.User) error {
	_, err := db.Exec("UPDATE users SET curAllocation=$1, targetAllocation=$2, taxRate=$3 WHERE email=$4", user.CurAllocation, user.TargetAllocation, user.TaxRate, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
