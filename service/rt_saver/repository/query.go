package repository

import (
	"context"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"strings"
)

func upsertTrades(ctx context.Context, tx pgx.Tx, trades []poloniex.RecentTrade) error {
	const (
		query = `
		INSERT INTO trades
			(tid, pair_id, price, amount, side_id, ts) VALUES %s
		ON CONFLICT DO NOTHING;`

		valuesTemplate = `(?,?,?,?,?,?)`
	)

	if len(trades) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(trades)*6)
	values := make([]string, 0, len(trades))
	for i := range trades {
		values = append(values, valuesTemplate)
		args = append(args,
			trades[i].Tid,
			poloniex.PairsToTypeMap[trades[i].Pair],
			trades[i].Price,
			trades[i].Amount,
			poloniex.SideToTypeMap[trades[i].Side],
			trades[i].Timestamp,
		)
	}

	valuesQuery := sqlx.Rebind(sqlx.DOLLAR, fmt.Sprintf(query, strings.Join(values, ",")))
	_, err := tx.Exec(ctx, valuesQuery, args...)
	return err
}
