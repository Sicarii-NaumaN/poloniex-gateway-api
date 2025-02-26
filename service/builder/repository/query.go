package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

func selectTrades(ctx context.Context, tx pgx.Tx, from, to int64) (resp []poloniex.RecentTrade, err error) {
	const (
		query = `
		SELECT 
			tid,
			pair_id,
			amount,
			price,
			side_id,
			ts
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
			&trade.Price,
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

func selectTradesByInterval(ctx context.Context, tx pgx.Tx, from, to int64, interval poloniex.Interval) (resp []poloniex.RecentTrade, err error) {
	// Сделано для ускорения разработки
	const (
		queryRaw = `
			SELECT 
				tid,
				pair_id,
				amount,
				price,
				side_id,
				ts
			FROM trades
			WHERE ts >= $1 AND ts < $2`

		queryOneMin     = ` AND NOT is_1m_processed`
		queryFifteenMin = ` AND NOT is_15m_processed`
		queryOneHour    = ` AND NOT is_1h_processed`
		queryOneDay     = ` AND NOT is_15d_processed`
		queryOrderBy    = ` ORDER BY ts;`
	)
	query := queryRaw
	switch interval {
	case poloniex.OneMin:
		query += queryOneMin + queryOrderBy
	case poloniex.FifteenMin:
		query += queryFifteenMin + queryOrderBy
	case poloniex.OneHour:
		query += queryOneHour + queryOrderBy
	case poloniex.OneDay:
		query += queryOneDay + queryOrderBy
	}

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
			&trade.Price,
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

func updateTradesByInterval(ctx context.Context, tx pgx.Tx, tids []string, interval poloniex.Interval) (err error) {
	// Сделано для ускорения разработки
	const (
		queryRaw = `
			UPDATE trades
			SET `

		queryOneMin     = ` is_1m_processed = true`
		queryFifteenMin = ` is_15m_processed = true`
		queryOneHour    = ` is_1h_processed = true`
		queryOneDay     = ` is_15d_processed = true`
		queryWhere      = ` WHERE tid = ANY($1)`
	)

	query := queryRaw
	switch interval {
	case poloniex.OneMin:
		query += queryOneMin + queryWhere
	case poloniex.FifteenMin:
		query += queryFifteenMin + queryWhere
	case poloniex.OneHour:
		query += queryOneHour + queryWhere
	case poloniex.OneDay:
		query += queryOneDay + queryWhere
	}

	_, err = tx.Exec(ctx, query, pq.Array(tids))
	if err != nil {
		return
	}

	return err
}

func selectLastSyncedCandleBegin(ctx context.Context, tx pgx.Tx) (lastSynced map[poloniex.Interval]int64, err error) {
	const query = `
			SELECT
			time_frame, MAX(begin_ts) AS max_begin_ts
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
		ON CONFLICT(pair_id, time_frame, begin_ts) DO UPDATE 
		    SET 
		        data = EXCLUDED.data;`

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
