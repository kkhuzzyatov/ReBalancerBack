package utils

import (
	"context"
	"errors"
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

func MoneyValueToFloat64(mv *investapi.MoneyValue) float64 {
	return float64(mv.Units) + float64(mv.Nano) / 1_000_000_000.0
}

func GetPrice(figi string, priceMap map[string]float64) (float64, error) {
	price, exists := priceMap[figi];
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
		// fmt.Println(err)
	}
	defer conn.Close()
	client := investapi.NewMarketDataServiceClient(conn)

	i := 0
	for i < 3 {
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
			return 0, err
		}
		if len(candlesResp.Candles) > 0 {
			return float64(candlesResp.Candles[0].Close.Units) + float64(candlesResp.Candles[0].Close.Nano)/1_000_000_000.0, nil
		}
	}
	
	return 0, errors.New("по позиции отсутствуют сделки за последние 3 дней")
}

func AllocationParser[num int | float64](alloc string) map[string]num {
	allocMap := make(map[string]num)
	for _, item := range strings.Split(alloc, ";") {
		parts := strings.Split(item, "=")
		var val num
		switch any(val).(type) {
		case int:
			val, _ := strconv.Atoi(parts[1])
			allocMap[parts[0]] = num(val)
		case float64:
			val, _ := strconv.ParseFloat(parts[1], 64)
			allocMap[parts[0]] = num(val)
		}
	}

	return allocMap
}

func СonvertKeysToUpperCase[V any](m map[string]V) map[string]V {
	result := make(map[string]V)
	for key, value := range m {
		result[strings.ToUpper(key)] = value
	}
	return result
}