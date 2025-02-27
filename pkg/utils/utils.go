package utils

import (
	"context"
	"errors"
	"fmt"
	"gomod/config"
	"strconv"
	"strings"
	"time"

	"github.com/ssummers02/invest-api-go-sdk/pkg/investapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertQuotationToFloat(q *investapi.Quotation) (float64, error) {
	if q == nil {
		return 0, errors.New("не удалось получить цену последней сделки")
	}
	return float64(q.Units) + float64(q.Nano)/1_000_000_000.0, nil
}

func CreatePriceMap(lastPrices []*investapi.LastPrice) map[string]float64 {
	priceMap := make(map[string]float64)
	for _, p := range lastPrices {
		if p.Price != nil {
			price, err := ConvertQuotationToFloat(p.Price)
			if err == nil {
				priceMap[p.Figi] = price
			}
		}
	}
	return priceMap
}

func GetPrice(figi string, priceMap map[string]float64) (float64, error) {
	price, exists := priceMap[figi]
	if exists && price != 0 {
		return price, nil
	}
	price, err := GetPriceOfCloseLastTradeDay(figi)
	if err != nil || price == 0 {
		return 0, err
	}
	return price, nil
}

func GetPriceOfCloseLastTradeDay(figi string) (float64, error) {
	conn, err := grpc.Dial("invest-public-api.tinkoff.ru:443", grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := investapi.NewMarketDataServiceClient(conn)

	i := 0
	for i < 5 {
		i++
		endTime := time.Now().AddDate(0, 0, 0)
		startTime := time.Now().AddDate(0, 0, -i)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer " + config.GetTokenOfTBank()))

		candlesResp, err := client.GetCandles(ctx, &investapi.GetCandlesRequest{
			Figi:     figi,
			From:     timestamppb.New(startTime),
			To:       timestamppb.New(endTime),
			Interval: investapi.CandleInterval_CANDLE_INTERVAL_DAY,
		})
		if err != nil {
			fmt.Println(err)
		}
		if len(candlesResp.Candles) != 0 {
			return float64(candlesResp.Candles[0].Close.Units) + float64(candlesResp.Candles[0].Close.Nano)/1_000_000_000.0, nil
		}
	}

	return 0, errors.New("по позиции отсутствуют сделки за последние 5 дней")
}

// Проверка удовлетворения налоговой ставки шаблону (xx,xx)
func CheckTaxRate(rate float64) bool {
	if 0 < rate && rate < 100 {
		return hasAtMostTwoDecimalPlaces(rate)
	}

	return false
}

// Проверка действительного числа на количнство знаков после запятой
func hasAtMostTwoDecimalPlaces(f float64) bool {
	formatted := strconv.FormatFloat(f, 'f', -1, 64)
	parts := strings.Split(formatted, ".")
	if len(parts) == 2 {
		return len(parts[1]) <= 2
	}

	return true
}
