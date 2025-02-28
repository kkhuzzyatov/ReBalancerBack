package repository

import (
	"database/sql"
	"gomod/config"
	"gomod/pkg/entities"
	"log"
)

var db *sql.DB

func InitDB() {
	var err error

	db, err = sql.Open("postgres", config.GetDBConnStr())
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to reach the database: %v", err)
	}
}

func CreateTable() {
	createTableStr := `CREATE TABLE IF NOT EXISTS users (
    id 							 	SERIAL PRIMARY KEY,
		email            	VARCHAR(255),
		password         	VARCHAR(255),
		curAllocation    	TEXT,
		targetAllocation 	TEXT,
		taxRate          	INT
	);`
	_, err := db.Exec(createTableStr)
	if err != nil {
			log.Fatalf("Error creating table: %v", err)
	}
}

func SelectUser(email string) (entities.User, error) {
	var user entities.User
	err := db.QueryRow("SELECT password, curAllocation, targetAllocation, taxRate FROM users WHERE email=$1", email).Scan(&user.Email, &user.PasswordHash, &user.CurAllocation, &user.TargetAllocation, &user.TaxRate)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func InsertUser(user entities.User) error {
	var id int
	err := db.QueryRow("INSERT INTO users (email, password, curAllocation, targetAllocation, taxRate) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Email, user.PasswordHash, user.CurAllocation, user.TargetAllocation, user.TaxRate).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(user entities.User) error {
	_, err := db.Exec("UPDATE users SET email=$2, password=$3, curAllocation=$4, targetAllocation=$5, taxRate=$6 WHERE id=$1", user.ID, user.Email, user.PasswordHash, user.CurAllocation, user.TargetAllocation, user.TaxRate)
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