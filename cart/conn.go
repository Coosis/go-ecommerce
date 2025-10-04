package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

func ConnectPoolWithBackoff(ctx context.Context, url string) (*pgxpool.Pool, error) {
	const (
		maxAttempts   = 10
		initialSleep  = 500 * time.Millisecond
		maxSleep      = 10 * time.Second
	)
	sleep := initialSleep

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		cfg, err := pgxpool.ParseConfig(url)
		if err != nil {
			return nil, err
		}

		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		pool, err := pgxpool.NewWithConfig(attemptCtx, cfg)
		cancel()

		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				return pool, nil
			} else {
				pool.Close()
				err = pingErr
			}
		}

		log.Warnf("DB not ready (attempt %d/%d): %v", attempt, maxAttempts, err)

		time.Sleep(sleep)
		sleep *= 2
		if sleep > maxSleep {
			sleep = maxSleep
		}
	}

	return nil, context.DeadlineExceeded
}

func MigrateWithBackoff(ctx context.Context, url string) (bool, error) {
	const (
		maxAttempts   = 10
		initialSleep  = 500 * time.Millisecond
		maxSleep      = 10 * time.Second
		perTryTO      = 5 * time.Second
	)

	db, err := sql.Open("pgx", url)
	if err != nil {
		return false, err
	}

	sleep := initialSleep
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, perTryTO)
		err = db.PingContext(attemptCtx)
		cancel()

		if err == nil {
			if err := goose.SetDialect("postgres"); err != nil {
				log.Fatal(err)
				return false, err
			}
			if err := goose.Up(db, "migrations"); err != nil {
				log.Fatal(err)
			}
			log.Info("migrations applied; starting API")
			return true, nil
		}

		log.Warnf("DB not ready (attempt %d/%d): %v", attempt, maxAttempts, err)

		select {
		case <-time.After(sleep):
			sleep *= 2
			if sleep > maxSleep {
				sleep = maxSleep
			}
		case <-ctx.Done():
			db.Close()
			return false, ctx.Err()
		}
	}

	db.Close()
	return false, context.DeadlineExceeded
}
