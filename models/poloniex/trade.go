package poloniex

import (
	"strconv"
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
	Timestamp int64  `json:"timestamp"` // время UTC UnixNano // Оставил милли
}
type RecentTradeArr []RecentTrade

func (trades RecentTradeArr) ConvertToKline(interval Interval) ([]Kline, error) {
	candles := make(map[int64]*Kline)
	//mu := new(sync.Mutex)

	for _, trade := range trades {
		price, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {

			return nil, err
		}
		amount, err := strconv.ParseFloat(trade.Amount, 64)
		if err != nil {
			return nil, err
		}

		startEnd := GetCandleIntervalByTime(trade.Timestamp, interval)

		startEnd.End-- // так возвращает сама биржа

		//mu.Lock()
		candle, exists := candles[startEnd.Start]
		if !exists {
			candle = &Kline{
				Pair:      trade.Pair,
				TimeFrame: interval,
				O:         price,
				H:         price,
				L:         price,
				C:         price,
				UtcBegin:  startEnd.Start,
				UtcEnd:    startEnd.End,
				VolumeBS:  VBS{},
			}
			candles[startEnd.Start] = candle
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
			candle.VolumeBS.BuyBase += amount / price
			candle.VolumeBS.BuyQuote += amount
		} else {
			candle.VolumeBS.SellBase += amount / price
			candle.VolumeBS.SellQuote += amount
		}
		//mu.Unlock()
	}

	var result []Kline
	for _, candle := range candles {
		result = append(result, *candle)
	}
	return result, nil
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
