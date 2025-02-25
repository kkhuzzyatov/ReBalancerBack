package repository

import (
	"context"
	"gomod/entities"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	BROKER_TYPE = "Т-Инвестиции"
)

type PGRepo struct {
	pool *pgxpool.Pool
}

func NewRepo(connString string) *PGRepo {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных")
	}

	return &PGRepo{
		pool: pool,
	}
}

// Регистрация нового пользователя
func (p *PGRepo) SignUp(user entities.User) error {
	queryString := "INSERT INTO users (id, tax_rate) VALUES ($1, $2);"
	_, err := p.pool.Exec(context.Background(), queryString, user.GetUserID(), user.GetTaxRate())
	if err != nil {
		return err
	}

	return nil
}

// Установление нового количества ценных бумаг по ticker и userID
func (p *PGRepo) SetCurAllocation(userID int64, allocation map[string]int) error {
	addQueryString := "INSERT INTO portfolio (ticker, amount, broker_type, user_id) VALUES ($1, $2, $3, $4);"
	updateQueryString := "UPDATE portfolio SET amount = $1 WHERE user_id = $2 AND ticker = $3;"

	for ticker, amount := range allocation {
		exists, err := p.CheckUserTicker(userID, ticker)
		if err != nil {
			return err
		}

		if exists {
			if _, err := p.pool.Exec(context.Background(), updateQueryString, amount, userID, ticker); err != nil {
				return err
			}
		} else {
			if _, err := p.pool.Exec(context.Background(), addQueryString, ticker, amount, BROKER_TYPE, userID); err != nil {
				return err
			}
		}
	}

	return nil
}

// Обновление налоговой ставки пользователя
func (p *PGRepo) SetTaxRate(userID int64, rate float64) error {
	queryString := "UPDATE users SET tax_rate = $1 WHERE id = $2;"
	_, err := p.pool.Exec(context.Background(), queryString, rate, userID)
	if err != nil {
		return err
	}

	return nil
}

// Проверка существования пользователя по userID
func (p *PGRepo) CheckID(userID int64) (bool, error) {
	var exists bool

	queryString := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);"
	err := p.pool.QueryRow(context.Background(), queryString, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Проверка существования тикера по userID
func (p *PGRepo) CheckUserTicker(userID int64, ticker string) (bool, error) {
	var exists bool

	queryString := "SELECT EXISTS(SELECT 1 FROM portfolio WHERE user_id = $1 AND ticker = $2);"
	err := p.pool.QueryRow(context.Background(), queryString, userID, ticker).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
