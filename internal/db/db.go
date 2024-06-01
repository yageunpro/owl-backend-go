package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

const TimeOut = 5 * time.Second

func Connect(dsn string) (*Pool, error) {
	pool, err := connectPool(context.Background(), dsn)
	if err != nil {
		return nil, errors.Join(ErrDBOpen, err)
	}
	if err := checkConnection(pool, TimeOut); err != nil {
		return nil, errors.Join(ErrDBPing, err)
	}
	return newPool(pool, nil), nil
}

func ConnectCluster(masterDsn, slaveDsn string) (*Pool, error) {
	masterPool, err := connectPool(context.Background(), masterDsn)
	if err != nil {
		return nil, errors.Join(ErrDBOpen, err)
	}
	if err := checkConnection(masterPool, TimeOut); err != nil {
		return nil, errors.Join(ErrDBPing, err)
	}
	slavePool, err := connectPool(context.Background(), slaveDsn)
	if err != nil {
		return nil, errors.Join(ErrDBOpen, err)
	}
	if err := checkConnection(slavePool, TimeOut); err != nil {
		return nil, errors.Join(ErrDBPing, err)
	}

	return newPool(masterPool, slavePool), nil
}

func connectPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	return pgxpool.NewWithConfig(ctx, dbConfig)
}

func checkConnection(pool *pgxpool.Pool, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return pool.Ping(ctx)
}
