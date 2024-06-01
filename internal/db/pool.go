package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	masterPgx *pgxpool.Pool
	slavePgx  *pgxpool.Pool
}

func newPool(master, slave *pgxpool.Pool) *Pool {
	if slave == nil {
		return &Pool{
			masterPgx: master,
			slavePgx:  nil,
		}
	}

	return &Pool{
		masterPgx: master,
		slavePgx:  slave,
	}
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.masterPgx.Exec(ctx, sql, args...)
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if p.slavePgx == nil {
		return p.masterPgx.Query(ctx, sql, args...)
	}
	return p.slavePgx.Query(ctx, sql, args...)
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if p.slavePgx == nil {
		return p.masterPgx.QueryRow(ctx, sql, args...)
	}
	return p.slavePgx.QueryRow(ctx, sql, args...)
}

func (p *Pool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.masterPgx.Begin(ctx)
}

func (p *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.masterPgx.BeginTx(ctx, txOptions)
}
