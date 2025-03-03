package db

import "time"

const (
	RetryMaxRetries int = 3
	RetryInitDelay      = 1 * time.Second
	RetryDeltaDelay     = 2 * time.Second
)

const driver string = "pgx"

const Timeout time.Duration = 1 * time.Second
