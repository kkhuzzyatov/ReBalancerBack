package handlers

import (
	"encoding/json"
	"fmt"
	"gomod/pkg/entities"
	"gomod/pkg/tBankAPI"
	"gomod/pkg/utils"
	"math"
	"strings"

	"net/http"
)

func CalcStd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user entities.CalcRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	tBankAPI.Mutex.Lock()
	response := CalcRebalance(user.CurAllocation, user.TargetAllocation)
	tBankAPI.Mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CalcRebalance(curAllocation []entities.CurAllocation, targetAlloc []entities.TargetAllocation) entities.CalcResponse {
	stocks := tBankAPI.Stocks

	var curAlloc = make(map[string]int)
	var targetAllocPercent = make(map[string]float64)
	for _, field := range curAllocation {
		curAlloc[field.Ticker] = field.Number
	}
	for _, field := range targetAlloc {
		targetAllocPercent[field.Ticker] = field.Percent
	}

	var response entities.CalcResponse
	curAlloc = utils.СonvertKeysToUpperCase(curAlloc)
	targetAllocPercent = utils.СonvertKeysToUpperCase(targetAllocPercent)

	_, exists := targetAllocPercent["RUB"]
	if !exists {
		targetAllocPercent["RUB"] = 0
	}

	var sumOfPercentInTargetAlloc float64 = 0
	for key, percent := range targetAllocPercent {
		_, exists := curAlloc[key]
		if !exists {
			curAlloc[key] = 0
		}
		sumOfPercentInTargetAlloc += percent
	}

	if 99.999 > sumOfPercentInTargetAlloc || sumOfPercentInTargetAlloc > 100.001 {
		response.Err = fmt.Sprintf("Ошибка: сумма процентов целевого распределения должна быть 100, а не %.2f.", sumOfPercentInTargetAlloc)
		return response
	}

	var totalBalance float64 = 0
	minLotPrice := math.MaxFloat64
	for ticker, number := range curAlloc {
		if stocks[ticker].Price <= 0 {
			response.Err = fmt.Sprintf("Ошибка: Не удалось получить данные актива с тикером %s.", ticker)
			return response
		}

		if stocks[ticker].Lot == 0 {
			continue
		}

		response.StockData = append(response.StockData, stocks[ticker])
		totalBalance += float64(number) * stocks[ticker].Price

		if minLotPrice > stocks[ticker].Price*float64(stocks[ticker].Lot) {
			minLotPrice = stocks[ticker].Price * float64(stocks[ticker].Lot)
		}
	}
	response.TotalValue = totalBalance
	cash := totalBalance

	targetAllocAmount := make(map[string]int)
	sellOrders := make(map[string]int)
	for ticker, number := range curAlloc {
		if stocks[ticker].Lot == 0 {
			continue
		}
		targetAllocAmount[ticker] = int(math.Floor(totalBalance * targetAllocPercent[ticker] * 0.01 / stocks[ticker].Price))
		sellOrders[ticker] = number - targetAllocAmount[ticker]
		if sellOrders[ticker] < 0 {
			sellOrders[ticker] = int(math.Ceil(float64(sellOrders[ticker]) / float64(stocks[ticker].Lot)))
		} else if sellOrders[ticker] != 0 {
			sellOrders[ticker] = int(math.Ceil(float64(sellOrders[ticker]) / float64(stocks[ticker].Lot)))
		}
		cash -= float64(number-sellOrders[ticker]*stocks[ticker].Lot) * stocks[ticker].Price
	}

	for cash > minLotPrice {
		for ticker := range targetAllocPercent {
			if stocks[ticker].Price*float64(stocks[ticker].Lot) < cash {
				sellOrders[ticker]--
				cash -= stocks[ticker].Price * float64(stocks[ticker].Lot)
			}
		}
	}

	for ticker := range curAlloc {
		if strings.ToUpper(ticker) == "RUB" {
			continue
		}
		if sellOrders[ticker] < 0 {
			response.Orders = append(response.Orders, entities.Order{
				Ticker:    ticker,
				NumberLot: -sellOrders[ticker],
			})
		} else if sellOrders[ticker] != 0 {
			response.Orders = append(response.Orders, entities.Order{
				Ticker:    ticker,
				NumberLot: sellOrders[ticker],
			})
		}
	}

	for ticker, number := range curAlloc {
		if stocks[ticker].Lot == 0 {
			continue
		}
		if number-sellOrders[ticker]*stocks[ticker].Lot > 0 {
			var numberOfAssets float64
			if ticker == "RUB" {
				numberOfAssets = float64(number-sellOrders[ticker]*stocks[ticker].Lot) + cash
			} else {
				numberOfAssets = float64(number - sellOrders[ticker]*stocks[ticker].Lot)
			}

			percentOfCapital := roundToHundredths(numberOfAssets * stocks[ticker].Price / totalBalance * 100)
			if ticker == "RUB" {
				response.FinalStructure = append(response.FinalStructure, entities.Position{
					Ticker:           ticker,
					Number:           numberOfAssets,
					PercentOfCapital: percentOfCapital,
				})
			} else {
				response.FinalStructure = append(response.FinalStructure, entities.Position{
					Ticker:           ticker,
					Number:           numberOfAssets,
					PercentOfCapital: percentOfCapital,
				})
			}
		}
	}

	return response
}

func roundToHundredths(num float64) float64 {
	return math.Round(num*100) / 100
}
