package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"strings"
)

func selectTrades(ctx context.Context, tx pgx.Tx, from, to int64) (resp []poloniex.RecentTrade, err error) {
	const (
		query = `
		SELECT 
			tid,
			pair_id,
			amount,
			side_id,
			ts,
		WHERE ts >= $1 
			AND ts < $2
			AND processed_at NOT NULL
		ORDER BY ts;`
	)

	rows, err := tx.Query(ctx, query, from, to)
	if err != nil {
		return
	}

	resp = make([]poloniex.RecentTrade, 0, 0)
	for rows.Next() {
		var (
			trade    poloniex.RecentTrade
			pairType poloniex.PairType
			sideType poloniex.SideType
		)

		err = rows.Scan(
			&trade.Tid,
			&pairType,
			&trade.Amount,
			&sideType,
			&trade.Timestamp,
		)
		if err != nil {
			return
		}

		trade.Pair = poloniex.TypePairsToMap[pairType]
		trade.Side = poloniex.TypeToSideMap[sideType]
		resp = append(resp, trade)
	}
	return resp, rows.Err()
}

func selectLastSyncedCandleBegin(ctx context.Context, tx pgx.Tx) (lastSynced map[poloniex.Interval]int64, err error) {
	const query = `
		SELECT time_frame, MAX(begin_ts) AS max_begin_ts
		FROM candles
		GROUP BY time_frame;
	`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return
	}

	lastSynced = make(map[poloniex.Interval]int64, 0)
	for rows.Next() {
		var (
			intervalType poloniex.IntervalType
			maxBeginTime int64
		)
		err = rows.Scan(&intervalType, &maxBeginTime)
		if err != nil {
			return
		}

		lastSynced[poloniex.TypeToInterval[intervalType]] = maxBeginTime
	}
	return lastSynced, rows.Err()
}

func upsertCandles(ctx context.Context, tx pgx.Tx, candles []poloniex.Kline) error {
	const (
		query = `
		INSERT INTO candles
			(pair_id, time_frame, begin_ts, end_ts, data) VALUES %s
		ON CONFLICT DO NOTHING;`

		valuesTemplate = `(?,?,?,?,?)`
	)

	if len(candles) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(candles)*5)
	values := make([]string, 0, len(candles))
	for i := range candles {
		data, err := json.Marshal(candles[i])
		if err != nil {
			return err
		}
		values = append(values, valuesTemplate)
		args = append(args,
			poloniex.PairsToTypeMap[candles[i].Pair],
			poloniex.IntervalToType[poloniex.Interval(candles[i].TimeFrame)],
			candles[i].UtcBegin,
			candles[i].UtcEnd,
			data,
		)
	}

	valuesQuery := sqlx.Rebind(sqlx.DOLLAR, fmt.Sprintf(query, strings.Join(values, ",")))
	_, err := tx.Exec(ctx, valuesQuery, args...)
	return err
}
