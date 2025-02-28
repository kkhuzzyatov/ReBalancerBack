package handlers

import "net/http"

func Handle() {
  http.HandleFunc("/save", corsMiddleware(Save))
  http.HandleFunc("/getdata", corsMiddleware(GetData))
  http.HandleFunc("/calc", corsMiddleware(Calc))
} 

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    next(w, r)
  }
}


