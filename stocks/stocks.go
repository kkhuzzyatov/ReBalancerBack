package stocks

import (
	"github.com/ssummers02/invest-api-go-sdk/pkg/investapi"
)

var (
	shares 	[]*investapi.Share
	bonds  	[]*investapi.Bond
	etfs 		[]*investapi.Etf

	priceMap map[string]float64
)

func SetShares(newShares []*investapi.Share) {
	shares = newShares
}

func SetBonds(newBonds []*investapi.Bond) {
	bonds = newBonds
}

func SetETFs(newETFs []*investapi.Etf) {
	etfs = newETFs
}

func SetPriceMap(newPriceMap map[string]float64) {
	priceMap = newPriceMap
}

func GetStockPriceByTicker(ticker string) float64 {
	return 0 //not done
}