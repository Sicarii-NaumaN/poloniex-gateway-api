package poloniex_api

import (
	"encoding/json"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"strconv"
)

// Тут все упрощено в рамках тестового задания
func repackGetCandleSticks(pair poloniex.Pair, respBody []byte) (poloniex.GetCandleStickResp, error) {
	var rawData [][]interface{}
	err := json.Unmarshal(respBody, &rawData)
	if err != nil {
		fmt.Println("-----------", string(respBody), "-----------")
		return poloniex.GetCandleStickResp{}, fmt.Errorf("error in json.Unmarshal: %w for pair: %s", err, string(pair))
	}

	var parsedData = make([]candleData, 0, len(rawData))
	for _, entry := range rawData {
		if len(entry) != 14 { // TODO more safe
			continue
		}

		dataEntry := candleData{
			Low:              entry[0].(string),
			High:             entry[1].(string),
			Open:             entry[2].(string),
			Close:            entry[3].(string),
			Amount:           entry[4].(string),
			Quantity:         entry[5].(string),
			BuyTakerAmount:   entry[6].(string),
			BuyTakerQuantity: entry[7].(string),
			TradeCount:       int(entry[8].(float64)),
			Timestamp:        int64(entry[9].(float64)),
			WeightedAverage:  entry[10].(string),
			Interval:         entry[11].(string),
			StartTime:        int64(entry[12].(float64)),
			CloseTime:        int64(entry[13].(float64)),
		}

		parsedData = append(parsedData, dataEntry)
	}
	return repackKL(pair, parsedData)
}

func repackKL(pair poloniex.Pair, in []candleData) (poloniex.GetCandleStickResp, error) {
	res := make([]poloniex.Kline, 0, len(in))
	for _, c := range in {
		o, err := strconv.ParseFloat(c.Open, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}

		cl, err := strconv.ParseFloat(c.Close, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}

		l, err := strconv.ParseFloat(c.Low, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}

		h, err := strconv.ParseFloat(c.High, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}

		buyTakerQuantity, err := strconv.ParseFloat(c.BuyTakerQuantity, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}
		buyTakerAmount, err := strconv.ParseFloat(c.BuyTakerAmount, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}
		quantity, err := strconv.ParseFloat(c.Quantity, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}
		amount, err := strconv.ParseFloat(c.Amount, 64)
		if err != nil {
			return poloniex.GetCandleStickResp{}, fmt.Errorf("error in strconv.ParseFloat: %w", err)
		}
		res = append(res, poloniex.Kline{
			Pair:      pair,
			TimeFrame: poloniex.Interval(c.Interval),
			O:         o,
			H:         h,
			L:         l,
			C:         cl,
			UtcBegin:  c.StartTime,
			UtcEnd:    c.CloseTime,
			VolumeBS: poloniex.VBS{
				BuyBase:   buyTakerQuantity,
				SellBase:  quantity - buyTakerQuantity,
				BuyQuote:  buyTakerAmount,
				SellQuote: amount - buyTakerAmount,
			},
		})
	}

	return poloniex.GetCandleStickResp{
		Candles: res,
		Total:   0,
	}, nil
}

func repackRT(pair poloniex.Pair, in []tradeData) (poloniex.GetTradesResp, error) {
	res := make([]poloniex.RecentTrade, 0, len(in))
	for _, t := range in {
		res = append(res, poloniex.RecentTrade{
			Tid:       t.Id,
			Pair:      pair,
			Price:     t.Price,
			Amount:    t.Amount,
			Side:      poloniex.Side(t.TakerSide),
			Timestamp: t.Ts,
		})
	}
	return poloniex.GetTradesResp{
		RecentTrades: res,
	}, nil
}
