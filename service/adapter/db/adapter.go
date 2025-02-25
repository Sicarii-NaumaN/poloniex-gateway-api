// Database adapter for sql server. Helps to make queries in on auth

package db

import (
	"context"
	"database/sql"
	"fmt"
	. "github.com/Sicarii-NaumaN/poloniex-gateway-api"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type IAdapter interface {
	InTx(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error
	GetConn(ctx context.Context) *pgxpool.Pool
}

type Adapter struct {
	Pool      *pgxpool.Pool
	Goose     *sql.DB
	isolation pgx.TxIsoLevel
}

func NewAdapter(gooseConn *sql.DB, conn *pgxpool.Pool, isolation pgx.TxIsoLevel) (IAdapter, error) {
	defer gooseConn.Close()

	ad := &Adapter{Pool: conn, Goose: gooseConn, isolation: isolation}
	err := ad.GooseUp()
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (b *Adapter) GetConn(ctx context.Context) *pgxpool.Pool {
	return b.Pool
}

func (b *Adapter) InTx(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) (err error) {
	tx, err := b.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: b.isolation})
	if err != nil {
		return fmt.Errorf("error creating tx: %s", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			logger.Error(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = f(ctx, tx)
	return
}

func (b *Adapter) GooseUp() error {
	goose.SetBaseFS(EmbedMigrations)
	if err := goose.Up(b.Goose, "migrations", goose.WithAllowMissing()); err != nil {
		return err
	}
	return nil
}

func (b *Adapter) GooseCreate() error {
	goose.SetBaseFS(EmbedMigrations)
	if err := goose.Create(b.Goose, "migrations", "", "sql"); err != nil {
		return err
	}
	return nil
}

func (b *Adapter) GooseDown() error {
	goose.SetBaseFS(EmbedMigrations)
	if err := goose.Down(b.Goose, "migrations"); err != nil {
		return err
	}

	return nil
}
