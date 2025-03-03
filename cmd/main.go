package main

import (
	"gomod/pkg/handlers"
	"gomod/pkg/repository"
	"gomod/pkg/stocks"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	stocks.InitializeData()

	repository.InitDB()

	http.HandleFunc("/calc", corsMiddleware(handlers.Calc))
  
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal(err)
  }
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}