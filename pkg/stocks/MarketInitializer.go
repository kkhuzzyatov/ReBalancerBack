package stocks

import (
	"fmt"
	"gomod/config"
	"gomod/pkg/entities"
	"gomod/pkg/utils"
	"log"

	"github.com/ssummers02/invest-api-go-sdk/pkg"
	"github.com/ssummers02/invest-api-go-sdk/pkg/investapi"
)

var Stocks map[string]entities.Stock
var PriceMap map[string]float64 

func InitializeData() {
	cfg := pkg.Config{
		Token:     config.GetTokenOfTBank(),
		AccountID: config.GetAccountID(),
	}

	srv, err := pkg.NewServicePool(cfg)
	if err != nil {
		log.Println(err)
	}

	Stocks = make(map[string]entities.Stock)

	shares, err := srv.GetSharesBase()
	if err != nil {
		log.Fatalf("Невозможно получить список акций, ошибка - %s", err)
	}
	CreateSharesMap(shares, Stocks)

	bonds, err := srv.GetBondsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список облигаций, ошибка - %s", err)
	}
	CreateBondsMap(bonds, Stocks)

	etfs, err := srv.GetETFsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список фондов, ошибка - %s", err)
	}
	CreateETFsMap(etfs, Stocks)
}

func CreateSharesMap(shares []*investapi.Share, Stocks map[string]entities.Stock) {
	for _, v := range shares {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag {
			price, err := utils.GetPrice(v.Figi, PriceMap)
			if err == nil {
				Stocks[v.Ticker] = entities.Stock{
					Figi: v.Figi,
					Lot: int(v.Lot),
					Price: price,
					Currency: v.Currency,
					AciValue: -1,
				}
			} else {
				fmt.Printf("Тикер: %s, Ошибка: %e\n", v.Ticker, err)
			}
		}
	}
}

func CreateBondsMap(bonds []*investapi.Bond, Stocks map[string]entities.Stock) {
	for _, v := range bonds {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag {
			price, err := utils.GetPrice(v.Figi, PriceMap)
			if err == nil {
				Stocks[v.Ticker] = entities.Stock{
					Figi: v.Figi,
					Lot: int(v.Lot),
					Price: price * 0.01 * utils.MoneyValueToFloat64(v.Nominal),
					Currency: v.Currency,
					AciValue: utils.MoneyValueToFloat64(v.AciValue),
				}
			} else {
				fmt.Printf("Тикер: %s, Ошибка: %e\n", v.Ticker, err)
			}
		}
	}
}

func CreateETFsMap(etfs []*investapi.Etf, Stocks map[string]entities.Stock) {
	for _, v := range etfs {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag {
			price, err := utils.GetPrice(v.Figi, PriceMap)
			if err == nil {
				Stocks[v.Ticker] = entities.Stock{
					Figi: v.Figi,
					Lot: int(v.Lot),
					Price: price,
					Currency: v.Currency,
					AciValue: -1,
				}
			} else {
				fmt.Printf("Тикер: %s, Ошибка: %e\n", v.Ticker, err)
			}
		}
	}
}

func IsAllocValid[num int | float64](alloc map[string]num) error {
	for ticker := range alloc {
		_, exists := Stocks[ticker]
		if !exists {
			return fmt.Errorf("не удалась получить данные о активе с тикером %s", ticker)
		}
	}
	return nil
}