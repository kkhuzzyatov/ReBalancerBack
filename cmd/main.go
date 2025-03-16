package main

import (
	"fmt"
	"gomod/pkg/handlers"
	"gomod/pkg/tBankAPI"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("программа запущена")
  
  fmt.Println("stocks loading is started")
	stockList := tBankAPI.FetchAssetsData()
	fmt.Println("stocks loading is done")
	tBankAPI.Stocks = stockList

	http.HandleFunc("/calc", handlers.Calc)
	
	go runServer()

	for {
		stockList := tBankAPI.FetchAssetsData()
		tBankAPI.Mutex.Lock()
		tBankAPI.Stocks = stockList
		tBankAPI.Mutex.Unlock()
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + ": updating stocks is done")
	}
}

func runServer(){
	fmt.Println("server is running")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}