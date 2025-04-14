package database

import (
	"fmt"
	"gomod/pkg/entities"
)

func CreateStocksTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS stocks (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			ticker VARCHAR(10) NOT NULL,
			number INTEGER NOT NULL,
			target_share FLOAT NOT NULL
		)
	`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы stocks: %v", err)
	}
	return nil
}

func InsertStock(stock entities.StockDB) error {
	query := `INSERT INTO stocks (user_id, ticker, number, target_share) VALUES ($1, $2, $3, $4)`
	_, err := DB.Exec(query, stock.UserID, stock.Ticker, stock.Number, stock.TargetShare)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении актива: %v", err)
	}
	return nil
}

func GetStocksByUserID(userID int) ([]entities.StockDB, error) {
	query := `SELECT id, user_id, ticker, number, target_share FROM stocks WHERE user_id = $1`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении актива: %v", err)
	}
	defer rows.Close()

	var stocks []entities.StockDB
	for rows.Next() {
		var stock entities.StockDB
		err := rows.Scan(&stock.ID, &stock.UserID, &stock.Ticker, &stock.Number, &stock.TargetShare)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании актива: %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func UpdateStock(stock entities.StockDB) error {
	query := `UPDATE stocks SET ticker = $1, number = $2, target_share = $3 WHERE id = $4 AND user_id = $5`
	_, err := DB.Exec(query, stock.Ticker, stock.Number, stock.TargetShare, stock.ID, stock.UserID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении актива: %v", err)
	}
	return nil
}

func DeleteStock(id, userID int) error {
	query := `DELETE FROM stocks WHERE id = $1 AND user_id = $2`
	_, err := DB.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении актива: %v", err)
	}
	return nil
}
