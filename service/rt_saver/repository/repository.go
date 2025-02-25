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

func (r *Repository) UpsertTrades(ctx context.Context, candles []poloniex.RecentTrade) (err error) {
	err = r.dbAdb.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = upsertTrades(ctx, tx, candles)
		if err != nil {
			return fmt.Errorf("error in upsertTrades: %w", err)
		}

		return nil
	})

	return
}
