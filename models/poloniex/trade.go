package poloniex

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	diffOneMin     = 60 * 1000
	diffFifteenMin = 15 * diffOneMin
	diffOneHour    = 4 * diffFifteenMin
	diffOneDay     = 24 * diffOneHour
)

type GetTradesReq struct {
	Pair  Pair
	Limit int64
}

type GetTradesResp struct {
	RecentTrades []RecentTrade
}

// Структура RT:
type RecentTrade struct {
	Tid       string `json:"tid"`       // id транзакции
	Pair      Pair   `json:"pair"`      // название валютной пары (как у нас)
	Price     string `json:"price"`     // цена транзакции
	Amount    string `json:"amount"`    // объём транзакции в базовой валюте
	Side      Side   `json:"side"`      // как биржа засчитала эту сделку (как buy или как sell)
	Timestamp int64  `json:"timestamp"` // время UTC UnixNano
}
type RecentTradeArr []RecentTrade

func (trades RecentTradeArr) AggregateToKlines(interval Interval) []Kline {
	candles := make(map[int64]*Kline)
	intervalMs := intervalDurationMs(interval)
	var mu sync.Mutex

	for _, trade := range trades {
		price := toFloat64(trade.Price)
		amount := toFloat64(trade.Amount)

		start := trade.Timestamp - (trade.Timestamp % intervalMs)
		end := start + intervalMs

		mu.Lock()
		candle, exists := candles[start]
		if !exists {
			candle = &Kline{
				Pair:      trade.Pair,
				TimeFrame: interval,
				O:         price,
				H:         price,
				L:         price,
				C:         price,
				UtcBegin:  start,
				UtcEnd:    end,
				VolumeBS:  VBS{},
			}
			candles[start] = candle
		}

		if price > candle.H {
			candle.H = price
		}
		if price < candle.L {
			candle.L = price
		}
		candle.C = price

		// Обновление объемов
		if trade.Side == Buy {
			candle.VolumeBS.BuyBase += amount
			candle.VolumeBS.BuyQuote += price * amount
		} else {
			candle.VolumeBS.SellBase += amount
			candle.VolumeBS.SellQuote += price * amount
		}
		mu.Unlock()
	}

	// Преобразование карты в массив
	var result []Kline
	for _, candle := range candles {
		result = append(result, *candle)
	}
	return result
}

func toFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования строки '%s' в float64: %v\n", s, err)
		return 0
	}
	return val
}

// intervalDurationMs возвращает длительность интервала в миллисекундах
func intervalDurationMs(interval Interval) int64 {
	switch IntervalToType[interval] {
	case IntervalTypeOneMin:
		return diffOneMin
	case IntervalTypeFifteenMin:
		return diffFifteenMin
	case IntervalTypeOneHour:
		return diffOneHour
	case IntervalTypeOneDay:
		return diffOneDay
	default:
		return diffOneMin
	}
}

//type Trade struct {
//	Amount0In    *big.Int `json:"amount0In"`
//	Amount1In    *big.Int `json:"amount1In"`
//	Amount0Out   *big.Int `json:"amount0Out"`
//	Amount1Out   *big.Int `json:"amount1Out"`
//	TxHash       string   `json:"tx_hash"`
//	CreationTime int64
//}
//
//type HistoryResponse struct {
//	History []TradeEvent `json:"history"`
//}
//
//type TradeEvent struct {
//	Date       int64   `json:"date"`
//	Type       string  `json:"type"`
//	TxHash     string  `json:"tx_hash"`
//	PriceUsd   float64 `json:"price_usd"`
//	TotalUsd   float64 `json:"total_usd"`
//	PriceQuote float64 `json:"price_quote"`
//	Spent      float64 `json:"spent"`
//	Received   float64 `json:"received"`
//}
