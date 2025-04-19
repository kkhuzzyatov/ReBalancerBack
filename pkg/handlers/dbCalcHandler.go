package handlers

import (
	"encoding/json"
	"gomod/pkg/database"
	"gomod/pkg/entities"
	"gomod/pkg/tBankAPI"
	"net/http"
)

func CalcDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var response entities.CalcResponse
	var userFromRequest entities.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userFromRequest)
	if err != nil {
		response.Err = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var userFromDB entities.UserDB
	userFromDB, err = database.GetUserByEmail(userFromRequest.Email)
	if err != nil {
		response.Err = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if userFromDB.Password != userFromRequest.Password {
		response.Err = "password is wrong"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var curAlloc []entities.CurAllocation
	var targetAlloc []entities.TargetAllocation

	stocks, _ := database.GetStocksByUserID(userFromDB.ID)
	for _, stock := range stocks {
		if stock.TargetShare == 0 {
			curAlloc = append(curAlloc, entities.CurAllocation{
				Ticker: stock.Ticker,
				Number: stock.Number})
		} else {
			targetAlloc = append(targetAlloc, entities.TargetAllocation{
				Ticker:  stock.Ticker,
				Percent: stock.TargetShare})
		}
	}

	tBankAPI.Mutex.Lock()
	response = CalcRebalance(curAlloc, targetAlloc)
	tBankAPI.Mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
