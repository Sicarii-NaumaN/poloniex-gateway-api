package repository

import (
	"context"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/adapter/db"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	dbAdb db.IAdapter
}

func NewRepository(dbAdb db.IAdapter) *Repository {
	return &Repository{dbAdb,
	}
}

func (r *Repository) SelectTradesByTS(ctx context.Context, from, to int64) (resp []poloniex.RecentTrade, err error) {
	err = r.dbAdb.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		resp, err = selectTrades(ctx, tx, from, to)
		if err != nil {
			return fmt.Errorf("error in selectLastSyncedCandle: %w", err)
		}

		return nil
	})

	return
}

func (r *Repository) SelectLastCandlesBeginTime(ctx context.Context) (lastSyncedTimeByInterval map[poloniex.Interval]int64, err error) {
	err = r.dbAdb.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		lastSyncedTimeByInterval, err = selectLastSyncedCandleBegin(ctx, tx)
		if err != nil {
			//if errors.Is(err, pgx.ErrNoRows) {
			//	return nil
			//}
			return fmt.Errorf("error in selectLastSyncedCandle: %w", err)
		}

		return nil
	})

	return
}

func (r *Repository) UpsertCandles(ctx context.Context, candles []poloniex.Kline) (err error) {
	err = r.dbAdb.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = upsertCandles(ctx, tx, candles)
		if err != nil {
			return fmt.Errorf("error in upsertCandles: %w", err)
		}

		return nil
	})

	return
}
