package main

import (
	"gomod/config"
	"gomod/pkg/stocks"
	"gomod/pkg/utils"
	"log"

	"github.com/ssummers02/invest-api-go-sdk/pkg"
)

func main() {
	cfg := pkg.Config{
		Token:     config.GetTokenOfTBank(),
		AccountID: config.GetAccountID(),
	}

	srv, err := pkg.NewServicePool(cfg)
	if err != nil {
		log.Println(err)
	}

	shares, err := srv.GetSharesBase()
	if err != nil {
		log.Fatalf("Невозможно получить список акций, ошибка - %s", err)
	}

	bonds, err := srv.GetBondsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список облигаций, ошибка - %s", err)
	}

	etfs, err := srv.GetETFsBase()
	if err != nil {
		log.Fatalf("Невозможно получить список фондов, ошибка - %s", err)
	}

	lastPrices, err := srv.GetLastPricesForAll()
	if err != nil {
		log.Fatalf("Невозможно получить последние цены, ошибка - %s", err)
	}
	priceMap := utils.CreatePriceMap(lastPrices)

	stocks.SetShares(shares)
	stocks.SetBonds(bonds)
	stocks.SetETFs(etfs)
	stocks.SetPriceMap(priceMap)

}