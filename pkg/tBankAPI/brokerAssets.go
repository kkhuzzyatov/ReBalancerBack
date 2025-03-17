package tBankAPI

import (
	"gomod/config"
	"gomod/pkg/entities"
	"gomod/pkg/utils"
	"log"
	"sync"

	"github.com/ssummers02/invest-api-go-sdk/pkg"
	"github.com/ssummers02/invest-api-go-sdk/pkg/investapi"
)

var PriceMap map[string]float64 
var Stocks map[string]entities.Stock
var Mutex sync.Mutex

func FetchAssetsData() map[string]entities.Stock {
	cfg := pkg.Config{
		Token:     config.GetTokenOfTBank(),
		AccountID: config.GetAccountID(),
	}

	srv, err := pkg.NewServicePool(cfg)
	if err != nil {
		log.Println(err)
	}

	stocks := make(map[string]entities.Stock)
	stocks["RUB"] = entities.Stock{
		Lot: 1, 
		Price: 1,
		AciValue: -1,
	}

	shares, err := srv.GetSharesBase()
	if err != nil {
		log.Fatalf("Невозможно получить список акций, ошибка - %s", err)
	}
	createSharesMap(shares, stocks)

	bonds, err := srv.GetBondsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список облигаций, ошибка - %s", err)
	}
	createBondsMap(bonds, stocks)

	etfs, err := srv.GetETFsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список фондов, ошибка - %s", err)
	}
	createETFsMap(etfs, stocks)

	return stocks
}

func createSharesMap(shares []*investapi.Share, stocks map[string]entities.Stock) {
	for _, v := range shares {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag && v.Currency == "rub" {
			price, err := utils.GetPrice(v.Figi, PriceMap)
			if err == nil && price != 0 {
				stocks[v.Ticker] = entities.Stock{
					Lot: int(v.Lot),
					Price: price,
					AciValue: -1,
				}
			} else {
				stocks[v.Ticker] = Stocks[v.Ticker]
			}
		}
	}
}

func createBondsMap(bonds []*investapi.Bond, stocks map[string]entities.Stock) {
	for _, v := range bonds {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag && v.Currency == "rub" {
			price, err := utils.GetPrice(v.Figi, PriceMap)
			if err == nil && price != 0 {
				stocks[v.Ticker] = entities.Stock{
					Lot: int(v.Lot),
					Price: price * 0.01 * utils.MoneyValueToFloat64(v.Nominal),
					AciValue: utils.MoneyValueToFloat64(v.AciValue),
				}
			} else {
				stocks[v.Ticker] = Stocks[v.Ticker]
			}
		}
	}
}

func createETFsMap(etfs []*investapi.Etf, stocks map[string]entities.Stock) {
	for _, v := range etfs {
		if v.ApiTradeAvailableFlag && v.BuyAvailableFlag && v.Currency == "rub" {
			price, err := utils.GetPrice(v.Figi, PriceMap)
			if err == nil && price != 0 {
				stocks[v.Ticker] = entities.Stock{
					Lot: int(v.Lot),
					Price: price,
					AciValue: -1,
				}
			} else {
				stocks[v.Ticker] = Stocks[v.Ticker]
			}
		}
	}
}