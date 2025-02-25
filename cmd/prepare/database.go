package prepare

import (
	"context"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/config"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/adapter/db"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/closer"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

func NewDBConn(ctx context.Context) db.IAdapter {
	dsn := config.GetConfigString(config.DBDSN)

	dbGoose, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		logger.Fatalf("can't open database goose driver adapter: %w", err.Error())
	}

	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatalf("can't open pgx database pool: %w", err.Error())

	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		logger.Fatalf("can't open pgx database pool: %w", err.Error())
	}

	closer.Add(pool)

	ad, err := db.NewAdapter(dbGoose, pool, pgx.ReadCommitted)
	if err != nil {
		logger.Fatal("can't initialize database adapter: %w", err.Error())
	}
	return ad
}
