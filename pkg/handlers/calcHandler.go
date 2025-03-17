package handlers

import (
	"encoding/json"
	"fmt"
	"gomod/pkg/entities"
	"gomod/pkg/tBankAPI"
	"gomod/pkg/utils"
	"math"

	"net/http"
)

func Calc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	
  if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
	}
	
	curAlloc := user.CurAllocation
	targetAlloc := user.TargetAllocation
	
	tBankAPI.Mutex.Lock()
	response := CalcRebalance(utils.AllocationParser[int](curAlloc), utils.AllocationParser[float64](targetAlloc), tBankAPI.Stocks)
	tBankAPI.Mutex.Unlock()

	type ResponseStruct struct {
		Response string `json:"response"`
	}

	responseStruct := ResponseStruct{
			Response: response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseStruct)
}

func CalcRebalance(curAlloc map[string]int, targetAllocPercent map[string]float64, stocks map[string]entities.Stock) string {
	curAlloc["RUB"] = 0
	result := ""
	curAlloc = utils.СonvertKeysToUpperCase(curAlloc)
	targetAllocPercent = utils.СonvertKeysToUpperCase(targetAllocPercent)

	var sumOfPercentInTargetAlloc float64 = 0
	for key, percent := range targetAllocPercent {
		_, exists := curAlloc[key]
		if !exists {
				curAlloc[key] = 0
		}
		sumOfPercentInTargetAlloc += percent
	}
	
	if sumOfPercentInTargetAlloc != 100 {
		result = fmt.Sprintf("Ошибка: сумма процентов целевого распределения должна быть 100, а не %.2f.", sumOfPercentInTargetAlloc)
		return result
	}

	result = "Актуальные данные по вашим активам:\n"
	var totalBalance float64 = 0
	minLotPrice := math.MaxFloat64
	for ticker, number := range curAlloc {
		if stocks[ticker].Price <= 0 {
			result = fmt.Sprintf("Ошибка: Не удалось получить данные актива с тикером %s.", ticker)
			return result
		}

		if stocks[ticker].AciValue != -1 {
			result += fmt.Sprintf("%s: Цена: %.2f, Накопленный купонный доход: %.2f, Полная стоимость: %.2f, Полная стоимость одного лота: %.2f\n", ticker, stocks[ticker].Price, stocks[ticker].AciValue, stocks[ticker].Price + stocks[ticker].AciValue, (stocks[ticker].Price + stocks[ticker].AciValue) * float64(stocks[ticker].Lot))
		} else {
			result += fmt.Sprintf("%s: Цена: %.2f, Полная стоимость одного лота: %.2f\n", ticker, stocks[ticker].Price, stocks[ticker].Price * float64(stocks[ticker].Lot))
		}
		
		if (stocks[ticker].Lot == 0) {
			continue
		}
		totalBalance += float64(number) * stocks[ticker].Price

		if minLotPrice > stocks[ticker].Price * float64(stocks[ticker].Lot) {
			minLotPrice = stocks[ticker].Price * float64(stocks[ticker].Lot)
		}
	}
	result += fmt.Sprintf("\nСовокупная стоимость активов с учётом пополнения: %.2f RUB\n", totalBalance)
	cash := totalBalance

	targetAllocAmount := make(map[string]int)
	sellOrders := make(map[string]int)
	for ticker, number := range curAlloc {
		if (stocks[ticker].Lot == 0) {
			continue
		}
		targetAllocAmount[ticker] = int(math.Floor(totalBalance * targetAllocPercent[ticker] * 0.01 / stocks[ticker].Price))
		sellOrders[ticker] = number - targetAllocAmount[ticker]
		if sellOrders[ticker] < 0 {
			sellOrders[ticker] = int(math.Ceil(float64(sellOrders[ticker]) / float64(stocks[ticker].Lot)))
		} else if sellOrders[ticker] != 0 {
			sellOrders[ticker] = int(math.Ceil(float64(sellOrders[ticker]) / float64(stocks[ticker].Lot)))
		}
		cash -= float64(number - sellOrders[ticker] * stocks[ticker].Lot) * stocks[ticker].Price
	}

	for cash > minLotPrice {
		for ticker := range targetAllocPercent {
			if stocks[ticker].Price * float64(stocks[ticker].Lot) < cash {
				sellOrders[ticker] --
				cash -= stocks[ticker].Price * float64(stocks[ticker].Lot)
			}
		}
	} 

	for ticker := range curAlloc {
		if sellOrders[ticker] < 0 {
			result += fmt.Sprintf("Купить %d лотов %s (%d штук)\n", -sellOrders[ticker], ticker, -sellOrders[ticker] * stocks[ticker].Lot)
		} else if sellOrders[ticker] != 0 {
			result += fmt.Sprintf("Продать %d лотов %s (%d штук)\n", sellOrders[ticker], ticker, sellOrders[ticker] * stocks[ticker].Lot)
		}
	}
	
	result += "\nИтоговая структура активов:\n"
	for ticker, number := range curAlloc {
		if (stocks[ticker].Lot == 0) {
			continue
		}
		if number - sellOrders[ticker] * stocks[ticker].Lot > 0 {
			var numberOfAssets float64
			if (ticker == "RUB") {
				numberOfAssets = float64(number - sellOrders[ticker] * stocks[ticker].Lot) + cash
			} else {
				numberOfAssets = float64(number - sellOrders[ticker] * stocks[ticker].Lot)
			}
			
			percentOfCapital := fmt.Sprintf("%.2f", roundToHundredths(numberOfAssets * stocks[ticker].Price / totalBalance * 100)) + "%"
			if (ticker == "RUB") {
				result += fmt.Sprintf("%s: %.2f (≈%s)\n", ticker, numberOfAssets, percentOfCapital)
			} else {
				result += fmt.Sprintf("%s: %d (≈%s)\n", ticker, int(numberOfAssets), percentOfCapital)
			}
		} 
	}

	return result
}

func roundToHundredths(num float64) float64 {
	return math.Round(num * 100) / 100
}