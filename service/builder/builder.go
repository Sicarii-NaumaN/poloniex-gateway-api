package builder

import (
	"context"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"time"
)

const (
	diffOneMin     = 60 * 1000
	diffFifteenMin = 15 * diffOneMin
	diffOneHour    = 4 * diffFifteenMin
	diffOneDay     = 24 * diffOneHour
)

type Repository interface {
	SelectLastCandlesBeginTime(ctx context.Context) (lastSyncedTimeByInterval map[poloniex.Interval]int64, err error)
	UpsertCandles(ctx context.Context, candles []poloniex.Kline) (err error)
	SelectTradesByInterval(ctx context.Context,
		candleIntervalsByTime map[poloniex.Interval]poloniex.StartEndInterval) (resp map[poloniex.Interval][]poloniex.RecentTrade, err error)
	UpdateTradesProcessedByInterval(ctx context.Context,
		tids []string, interval poloniex.Interval) (err error)
}

// Тут забираются трейды и преобразуются в свечи
type Service struct {
	rep Repository

	intervalDiffs map[poloniex.Interval]int64
}

func NewService(rep Repository) *Service {
	return &Service{
		rep: rep,
		intervalDiffs: map[poloniex.Interval]int64{
			poloniex.OneMin:     diffOneMin,
			poloniex.FifteenMin: diffFifteenMin,
			poloniex.OneHour:    diffOneHour,
			poloniex.OneDay:     diffOneDay,
		},
	}
}

const batchSize = 10000

func (s *Service) Build(ctx context.Context, candleIntervalsByTime map[poloniex.Interval]poloniex.StartEndInterval) error {
	tradesByInterval, err := s.rep.SelectTradesByInterval(ctx, candleIntervalsByTime)
	if err != nil {
		return fmt.Errorf("error in SelectTradesByInterval: %w", err)
	}

	// Знаю что так делать не лучшая затея, допущение в рамках тестового + знаем точное кол-во итераций
	for i, trades := range tradesByInterval {

		tradesByPair := splitTradesByPair(trades)

		for _, tradesPair := range tradesByPair {
			tradesArr := poloniex.RecentTradeArr(tradesPair)

			var candles []poloniex.Kline
			candles, err = tradesArr.ConvertToKline(i)
			if err != nil {
				return fmt.Errorf("error in tradesArr.ConvertToKline: %w", err)
			}

			err = s.rep.UpsertCandles(ctx, candles)
			if err != nil {
				return fmt.Errorf("error in UpsertCandles: %w", err)
			}

		}

		// Это значит, что мы все трейды в данной свечи уже учли и можем не процессить больше
		// Возможно даже излишне, замедляет запрос
		if time.Now().UTC().UnixMilli() > candleIntervalsByTime[i].End {
			tids := make([]string, 0, len(trades))
			for _, trade := range trades {
				tids = append(tids, trade.Tid)
			}

			//for j := 0; j < len(tids); j += batchSize {
			//	end := j + batchSize
			//	if end > len(tids) {
			//		end = len(tids)
			//	}
			//	batch := tids[j:end]
			err = s.rep.UpdateTradesProcessedByInterval(ctx, tids, i)
			if err != nil {
				return fmt.Errorf("error in UpdateTradesProcessedByInterval: %w", err)
			}
			//}
		}
	}

	return nil
}

func splitTradesByPair(trades []poloniex.RecentTrade) map[poloniex.Pair][]poloniex.RecentTrade {
	var resp = make(map[poloniex.Pair][]poloniex.RecentTrade, len(poloniex.TypePairsToMap))
	for _, trade := range trades {
		if tr, ok := resp[trade.Pair]; !ok {
			resp[trade.Pair] = []poloniex.RecentTrade{trade}
		} else {
			resp[trade.Pair] = append(tr, trade)
		}
	}
	return resp
}
