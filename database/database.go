package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	connString = "postgres://postgres:postgres@localhost:5432/dubmasterdb"
	Pool       *pgxpool.Pool
)

func ConnectDatabase() error {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("unable to ping database: %v", err)
	}

	Pool = pool

	slog.Info("[Dubmaster] Conected with database successfully")
	return nil
}

func CloseDatabase() {
	Pool.Close()
}

func GetConnection() *pgxpool.Conn {
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
