package utils

import (
	"context"
	"errors"
	"fmt"
	"gomod/config"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/ssummers02/invest-api-go-sdk/pkg/investapi"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func QuotationToFloat(q *investapi.Quotation) (float64, error) {
	if q == nil {
		return 0, errors.New("не удалось получить цену последней сделки")
	}
	return float64(q.Units) + float64(q.Nano)/1_000_000_000.0, nil
}

func MoneyValueToFloat64(mv *investapi.MoneyValue) float64 {
	return float64(mv.Units) + float64(mv.Nano) / 1_000_000_000.0
}

func CreatePriceMap(lastPrices []*investapi.LastPrice) map[string]float64 {
	priceMap := make(map[string]float64)
	for _, p := range lastPrices {
		if p.Price != nil {
			price, err := QuotationToFloat(p.Price)
			if err == nil {
				priceMap[p.Figi] = price
			}
		}
	}
	return priceMap
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
			return 0, err
		}
		if len(candlesResp.Candles) > 0 {
			return float64(candlesResp.Candles[0].Close.Units) + float64(candlesResp.Candles[0].Close.Nano)/1_000_000_000.0, nil
		}
	}
	
	return 0, errors.New("по позиции отсутствуют сделки за последние 5 дней")
}

func CheckEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CheckPassword(password string) bool {
	return len(password) >= 8
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashAndPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}