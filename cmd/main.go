package main

import (
	"gomod/pkg/handlers"
	"gomod/pkg/repository"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// stocks.InitializeData()

	repository.InitDB()

	handlers.Handle()
  
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal(err)
  }
}