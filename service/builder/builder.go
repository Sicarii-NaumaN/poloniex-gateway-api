package builder

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
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

func (s *Service) Build(ctx context.Context, from int64) error {

	return nil
}
