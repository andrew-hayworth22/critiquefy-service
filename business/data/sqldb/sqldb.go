package sqldb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

var (
	ErrDBNotFound        = pgx.ErrNoRows
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrUndefinedTable    = errors.New("undefined table")
)

// Config defines what is needed to connect to the database
type Config struct {
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

// Open creates a connection pool to the database based on the configuration
func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	fmt.Println(cfg.URL)
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB: %w", err)
	}
	return pool, nil
}

// StatusCheck checks the connection to the database
func StatusCheck(ctx context.Context, pool *pgxpool.Pool) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = pool.Ping(ctx)
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	const q = `SELECT true;`
	var tmp bool
	return pool.QueryRow(ctx, q).Scan(&tmp)
}
