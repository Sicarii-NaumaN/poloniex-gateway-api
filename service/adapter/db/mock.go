package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mock struct{}

func NewDBMock() IAdapter {
	return &mock{}
}

func (m *mock) GetConn(ctx context.Context) *pgxpool.Pool {
	return nil
}

func (m *mock) InTx(context.Context, func(ctx context.Context, tx pgx.Tx) error) error {
	return nil
}
