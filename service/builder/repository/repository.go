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

func (r *Repository) SelectTradesByInterval(ctx context.Context,
	candleIntervalsByTime map[poloniex.Interval]poloniex.StartEndInterval) (resp map[poloniex.Interval][]poloniex.RecentTrade, err error) {

	resp = make(map[poloniex.Interval][]poloniex.RecentTrade)
	err = r.dbAdb.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Знаю что так делать не лучшая затея, допущение в рамках тестового + знаем точное кол-во итераций
		for interval, candleInterval := range candleIntervalsByTime {
			resp[interval], err = selectTradesByInterval(ctx, tx, candleInterval.Start, candleInterval.End, interval)
			if err != nil {
				return fmt.Errorf("error in selectTradesByInterval with interval %s start %d end %d: %w",
					interval, candleInterval.Start, candleInterval.End, err)
			}
		}

		return nil
	})

	return
}

func (r *Repository) UpdateTradesProcessedByInterval(ctx context.Context,
	tids []string, interval poloniex.Interval) (err error) {

	err = r.dbAdb.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = updateTradesByInterval(ctx, tx, tids, interval)
		if err != nil {
			//if errors.Is(err, pgx.ErrNoRows) {
			//	return nil
			//}
			return fmt.Errorf("error in updateTradesByInterval: %w", err)
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
