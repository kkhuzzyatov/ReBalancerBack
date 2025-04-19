package main

import (
	"fmt"
	"gomod/pkg/database"
	"gomod/pkg/entities"
	"gomod/pkg/handlers"
	"gomod/pkg/tBankAPI"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("программа запущена")

	err := database.InitDB()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}

	err = database.CreateUsersTable()
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы пользователей: %v", err)
	}

	err = database.CreateStocksTable()
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы активов: %v", err)
	}

	fmt.Println("stocks loading is started")
	stockList := tBankAPI.FetchAssetsData()
	stockList["FOR_TEST_0"] = entities.Stock{
		Ticker:   "FOR_TEST_0",
		Lot:      1,
		Price:    6.05,
		AciValue: -1,
	}
	stockList["FOR_TEST_1"] = entities.Stock{
		Ticker:   "FOR_TEST_1",
		Lot:      10,
		Price:    325,
		AciValue: -1,
	}
	stockList["RUB"] = entities.Stock{
		Ticker:   "RUB",
		Lot:      1,
		Price:    1,
		AciValue: -1,
	}
	fmt.Println("stocks loading is done")
	tBankAPI.Stocks = stockList

	http.HandleFunc("/calc", handlers.CalcStd)
	http.HandleFunc("/signup", handlers.CreateUser)
	http.HandleFunc("/finduser", handlers.FindUser)
	http.HandleFunc("/deleteuser", handlers.DeleteUser)
	http.HandleFunc("/saveuserdata", handlers.SaveHandler)
	http.HandleFunc("/dbcalc", handlers.CalcDB)

	go runServer()

	for {
		time.Sleep(4 * time.Hour)
		stockList := tBankAPI.FetchAssetsData()
		stockList["FOR_TEST_0"] = entities.Stock{
			Ticker:   "FOR_TEST_0",
			Lot:      1,
			Price:    6.05,
			AciValue: -1,
		}
		stockList["FOR_TEST_1"] = entities.Stock{
			Ticker:   "FOR_TEST_1",
			Lot:      10,
			Price:    325,
			AciValue: -1,
		}
		stockList["RUB"] = entities.Stock{
			Ticker:   "RUB",
			Lot:      1,
			Price:    1,
			AciValue: -1,
		}
		tBankAPI.Mutex.Lock()
		tBankAPI.Stocks = stockList
		tBankAPI.Mutex.Unlock()
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + ": updating stocks is done")
	}
}

func runServer() {
	fmt.Println("server is running")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
