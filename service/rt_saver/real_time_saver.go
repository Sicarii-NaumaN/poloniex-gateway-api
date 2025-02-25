package rt_saver

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
)

const batchSize = 100

type Repository interface {
	UpsertTrades(ctx context.Context, candles []poloniex.RecentTrade) (err error)
}

type Service struct {
	rep Repository
}

func NewService(rep Repository) *Service {
	return &Service{rep: rep}
}

func (s *Service) SaveRealTimeTrades(ctx context.Context, tradeChan <-chan poloniex.RecentTrade) error {
	tradeBuffer := make([]poloniex.RecentTrade, 0, batchSize)

	for {
		select {
		case <-ctx.Done():
			// Сохраняем оставшиеся
			if len(tradeBuffer) > 0 {
				err := s.rep.UpsertTrades(ctx, tradeBuffer)
				if err != nil {
					return err
				}
				logger.Infof("inserting %d trades...", len(tradeBuffer))
			}
			return nil
		case trade := <-tradeChan:
			tradeBuffer = append(tradeBuffer, trade)

			// Кладем бачами
			if len(tradeBuffer) >= batchSize {
				err := s.rep.UpsertTrades(ctx, tradeBuffer)
				if err != nil {
					return err
				}
				logger.Info(ctx, "inserting 100 trades....")
				tradeBuffer = tradeBuffer[:0]
			}
		}
	}
}
