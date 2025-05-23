package sqldb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
	URL string
}

// Open creates a connection to the database based on the configuration
func Open(ctx context.Context, cfg Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB: %w", err)
	}
	return conn, nil
}

// StatusCheck checks the connection to the database
func StatusCheck(ctx context.Context, conn *pgx.Conn) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = conn.Ping(ctx)
		if pingError != nil {
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
	return conn.QueryRow(ctx, q).Scan(&tmp)

}
