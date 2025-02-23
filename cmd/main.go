package main

import (
	"gomod/config"
	"gomod/stocks"
	"gomod/utils"
	"log"

	"github.com/ssummers02/invest-api-go-sdk/pkg"
)

func main() {
	cfg := pkg.Config{
		Token:     config.GetToken(),
		AccountID: config.GetAccountID(),
	}

	srv, err := pkg.NewServicePool(cfg)
	if err != nil {
		log.Println(err)
	}

	lastPrices, err := srv.GetLastPricesForAll()
	if err != nil {
		log.Fatalf("Невозможно получить последние цены, ошибка - %s", err)
	}
	priceMap := utils.CreatePriceMap(lastPrices)

	shares, err := srv.GetSharesBase()
	if err != nil {
		log.Fatalf("Невозможно получить список акций, ошибка - %s", err)
	}
	stocks.PrintShares(shares, priceMap)

	etfs, err := srv.GetETFsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список фондов, ошибка - %s", err)
	}
	stocks.PrintETFs(etfs, priceMap)

	bonds, err := srv.GetBondsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список облигаций, ошибка - %s", err)
	}
	stocks.PrintBonds(bonds, priceMap)
}