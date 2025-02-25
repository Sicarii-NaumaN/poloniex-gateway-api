package processor

import (
	"context"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"math"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/errgroup"
	"sync"
	"time"
)

var initSyncDate = time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

const (
	limit = 500

	diffOneMin     = 60 * 1000 * limit
	diffFifteenMin = 15 * diffOneMin
	diffOneHour    = 4 * diffFifteenMin
	diffOneDay     = 24 * diffOneHour
)

type Repository interface {
	SelectLastCandlesBeginTime(ctx context.Context) (lastSyncedTimeByInterval map[poloniex.Interval]int64, err error)
	UpsertCandles(ctx context.Context, candles []poloniex.Kline) (err error)
}

type PoloniexAdapter interface {
	GetCandleSticks(ctx context.Context, req poloniex.GetCandleStickReq) (poloniex.GetCandleStickResp, error)
}

type Service struct {
	cexAdb PoloniexAdapter
	rep    Repository

	intervalDiffs map[poloniex.Interval]int64
}

func NewService(cexAdb PoloniexAdapter, rep Repository) *Service {
	return &Service{
		cexAdb: cexAdb,
		rep:    rep,
		intervalDiffs: map[poloniex.Interval]int64{
			poloniex.OneMin:     diffOneMin,
			poloniex.FifteenMin: diffFifteenMin,
			poloniex.OneHour:    diffOneHour,
			poloniex.OneDay:     diffOneDay,
		},
	}
}

// Process
// В данной реализации решил все делать в одном процессе
// Складировать в одной базе, тем самым не создавая конкурентное состояние при записи
func (s *Service) Process(ctx context.Context, finishChan chan struct{}, timeBeforeAll time.Time) error {
	// Достаем для каждого интервала время последней засинхронизированной свечи
	// По хорошему можно in-memory хранить, но не успел сделать
	lastSyncedCandleBeginTime, err := s.rep.SelectLastCandlesBeginTime(ctx)
	if err != nil {
		return fmt.Errorf("error in SelectLastCandle: %w", err)
	}

	// Если для какого-то интервала не было засинканых свечей, то проставляем время запуска приложения
	for k, _ := range s.intervalDiffs {
		if lastSyncedCandleBeginTime[k] == 0 {
			lastSyncedCandleBeginTime[k] = initSyncDate
		}
	}

	candles, err := s.makeRequests(ctx, lastSyncedCandleBeginTime)
	if err != nil {
		return fmt.Errorf("error in makeRequests: %w", err)
	}

	// Если для какого-то интервала не было засинканых свечей, то проставляем время запуска приложения
	for k, _ := range s.intervalDiffs {
		if lastSyncedCandleBeginTime[k] == 0 {
			lastSyncedCandleBeginTime[k] = initSyncDate
		}
	}

	err = s.rep.UpsertCandles(ctx, candles)
	if err != nil {
		return fmt.Errorf("error in UpsertCandles: %w", err)
	}

	// Если синк дошел до времени, когда запустилось приложение, то последующее схватит вебсокет
	if s.isFinishedSync(timeBeforeAll.UnixMilli(), lastSyncedCandleBeginTime) {
		go func() {
			finishChan <- struct{}{}
			close(finishChan)
		}()

	}
	return nil
}

// Можно покрыть юнитами
func (s *Service) makeRequests(ctx context.Context, intervalsLastSyncedCandle map[poloniex.Interval]int64) ([]poloniex.Kline, error) {
	var (
		eg = new(errgroup.Group)
		mu = new(sync.Mutex)
	)

	var candles = make([]poloniex.Kline, 0, limit*len(poloniex.PairsToTypeMap))
	for interval := range poloniex.IntervalToType {
		start := intervalsLastSyncedCandle[interval]
		end := start + s.intervalDiffs[interval]

		for pair := range poloniex.PairsToTypeMap {
			eg.Go(func() error {
				return s.fetchCandles(ctx, mu, pair, interval, start, end, &candles)
			})
		}

	}

	if err := eg.Wait(); err != nil {
		return candles, fmt.Errorf("error in fetchCandles: %w", err)
	}

	return candles, nil

}

// Можно покрыть юнитами
func (s *Service) fetchCandles(ctx context.Context, mu *sync.Mutex,
	pair poloniex.Pair, interval poloniex.Interval,
	start, end int64,
	dst *[]poloniex.Kline) error {
	resp, err := s.cexAdb.GetCandleSticks(ctx, poloniex.GetCandleStickReq{
		Pair:      pair,
		Interval:  interval,
		StartTime: start,
		EndTime:   end,
		Limit:     limit,
	})
	if err != nil {
		return fmt.Errorf("error in GetCandleSticks: %v", err)
	}

	logger.Infof("Fetched %d candles for pair: %s, interval: %s with start: %d, end: %d",
		len(resp.Candles), pair, interval, start, end)

	mu.Lock()
	defer mu.Unlock()

	*dst = append(*dst, resp.Candles...)
	return nil
}

// Можно покрыть юнитами
func (s *Service) isFinishedSync(timeBeforeAll int64,
	lastSyncedCandleBeginTime map[poloniex.Interval]int64) bool {

	// У каждой свечи смотрим, что текущая итерация была самой актуальной, сравнивая со временем запуска приложения
	// и останавливаем процессинг
	for k, v := range s.intervalDiffs {
		lastSyncedCandleBeginTime[k] += v
	}

	minSyncedCandleBeginTime := minSyncedTimeFrame(lastSyncedCandleBeginTime)
	if timeBeforeAll <= minSyncedCandleBeginTime {
		return true
	}
	return false
}

// Можно покрыть юнитами
func minSyncedTimeFrame(m map[poloniex.Interval]int64) int64 {
	minVal := int64(math.MaxInt64)
	for _, v := range m {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}
