package stocks

import (
	"fmt"
	"gomod/utils"

	"github.com/ssummers02/invest-api-go-sdk/pkg/investapi"
)

func PrintShares(shares []*investapi.Share, priceMap map[string]float64) {
	fmt.Println("=== Акции ===")
	for _, v := range shares {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag && v.SellAvailableFlag {
			price, err := utils.GetPrice(v.Figi, priceMap)
			if err == nil {
				fmt.Printf("Тикер: %s, Лотность: %d, Цена: %.2f, Валюта расчётов: %s\n",
				v.Ticker, v.Lot, price, v.Currency)
			} else {
				fmt.Printf("Тикер: %s, Ошибка: %e\n", v.Ticker, err)
			}
		}
	}
}

func PrintETFs(etfs []*investapi.Etf, priceMap map[string]float64) {
	fmt.Println("=== Фонды (ETF) ===")
	for _, v := range etfs {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag && v.SellAvailableFlag {
			price, err := utils.GetPrice(v.Figi, priceMap)
			if err == nil {
				fmt.Printf("Тикер: %s, Лотность: %d, Цена: %.2f, Валюта расчётов: %s\n",
				v.Ticker, v.Lot, price, v.Currency)
			} else {
				fmt.Printf("Тикер: %s, Ошибка: %e\n", v.Ticker, err)
			}
		}
	}
}

func PrintBonds(bonds []*investapi.Bond, priceMap map[string]float64) {
	fmt.Println("=== Облигации ===")
	for _, v := range bonds {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag && v.SellAvailableFlag {
			price, err := utils.GetPrice(v.Figi, priceMap)
			if err == nil {
				fmt.Printf("Тикер: %s, Лотность: %d, Цена: %.2f, Валюта расчётов: %s\n",
				v.Ticker, v.Lot, price, v.Currency)
			} else {
				fmt.Printf("Тикер: %s, Ошибка: %e\n", v.Ticker, err)
			}
		}
	}
}