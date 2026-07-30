package main

import "time"

type Job struct {
	ID         int
	Priority   int
	Duration   time.Duration
	Retries    int
	MaxRetries int
	Timeout    time.Duration
}
