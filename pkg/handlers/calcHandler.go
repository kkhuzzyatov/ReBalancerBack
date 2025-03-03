package handlers

import (
	"encoding/json"
	"fmt"
	"gomod/pkg/entities"
	"gomod/pkg/stocks"
	"gomod/pkg/utils"

	// "gomod/pkg/utils"
	"net/http"
)

func Calc(w http.ResponseWriter, r *http.Request) {
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

	// 	err = stocks.IsAllocValid[int](utils.AllocationParser[int](user.CurAllocation)) 
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// }
	// 	err = stocks.IsAllocValid[float64](utils.AllocationParser[float64](user.TargetAllocation)) 
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	response := CalcRebalance(utils.AllocationParser[int](curAlloc), utils.AllocationParser[float64](targetAlloc), stocks.Stocks)

	w.Write([]byte(response))
}

func CalcRebalance(curAlloc map[string]int, targetAllocPercent map[string]float64, Stocks map[string]entities.Stock) string {
	result := "Актуальные данные по текущий конфигурации активов портфеля:\n"
	var totalBalance float64 = 0
	for ticker, number := range curAlloc {
		if Stocks[ticker].AciValue != -1 {
			result += fmt.Sprintf("%s: Цена: %.2f, Накопленный купонный доход: %.2f, Полная стоимость: %.2f, Полная стоимость одного лота: %.2f\n", ticker, Stocks[ticker].Price, Stocks[ticker].AciValue, Stocks[ticker].Price + Stocks[ticker].AciValue, (Stocks[ticker].Price + Stocks[ticker].AciValue) * float64(Stocks[ticker].Lot))
		} else {
			result += fmt.Sprintf("%s: Цена: %.2f, Полная стоимость одного лота: %.2f\n", ticker, Stocks[ticker].Price, Stocks[ticker].Price * float64(Stocks[ticker].Lot))
		}
		
		if (Stocks[ticker].Lot == 0) {
			continue
		}
		totalBalance += float64(number) * Stocks[ticker].Price
	}
	result += fmt.Sprintf("\nСовокупная стоимость активов: %.2f RUB\n", totalBalance)

	result += "\nЧтобы привести портфель в соответствие с целевыми пропорциями нужно:\n"
	targetAllocAmount := make(map[string]int)
	sellOrders := make(map[string]int)
	for ticker, number := range curAlloc {
		if (Stocks[ticker].Lot == 0) {
			continue
		}
		targetAllocAmount[ticker] = int(totalBalance * targetAllocPercent[ticker] * 0.01 / Stocks[ticker].Price)
		sellOrders[ticker] = int((number - targetAllocAmount[ticker]) / Stocks[ticker].Lot)
		if sellOrders[ticker] < 0 {
			result += fmt.Sprintf("Купить %d лотов %s\n", -sellOrders[ticker], ticker)
		} else if sellOrders[ticker] != 0 {
			result += fmt.Sprintf("Продать %d лотов %s\n", sellOrders[ticker], ticker)
		}
	}

	result += "\nИтоговая структура активов:\n"
	for ticker, number := range curAlloc {
		if (Stocks[ticker].Lot == 0) {
			continue
		}
		result += fmt.Sprintf("%s: %d\n", ticker, number - sellOrders[ticker])
	}

	return result
}