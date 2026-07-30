package main

import (
	"context"
	"time"
)

func ExecuteJob(ctx context.Context, job Job) error {
	select {
	case <-time.After(job.Duration):
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
