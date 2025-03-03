package handlers

import "net/http"

func Handle() {
  http.HandleFunc("/save", corsMiddleware(Save))
  http.HandleFunc("/get", corsMiddleware(GetData))
  http.HandleFunc("/calc", corsMiddleware(Calc))
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


