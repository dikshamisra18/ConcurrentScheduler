package main

import (
	"math/rand"
	"time"
)

type BenchmarkRequest struct {
	Workers int `json:"workers"`
	Jobs    int `json:"jobs"`
}

type BenchmarkResult struct {
	Workers     int     `json:"workers"`
	Jobs        int     `json:"jobs"`
	Seconds     float64 `json:"seconds"`
	Throughput  float64 `json:"throughput"`
	Utilization float64 `json:"utilization"`
	Completed   int     `json:"completed"`
	Failed      int     `json:"failed"`
}

func StartBenchmark(workers int, jobs int) {

	// Stop existing scheduler
	if scheduler != nil {
		scheduler.Cancel()
		scheduler.WG.Wait()
	}

	// create fresh scheduler
	scheduler = NewScheduler()

	// record time
	start := time.Now()

	// start workers
	scheduler.StartWorkers(workers)

	// start dispatcher
	go scheduler.Dispatch()

	// submit jobs
	for i := 1; i <= jobs; i++ {
		job := Job{
			ID:         i,
			Priority:   rand.Intn(10) + 1,
			Duration:   time.Duration(rand.Intn(3)+1) * time.Second,
			Timeout:    time.Duration(rand.Intn(3)+2) * time.Second,
			Retries:    0,
			MaxRetries: 3,
		}
		scheduler.Submit(job)
	}

	for {
		scheduler.Mutex.Lock()

		done := scheduler.Metrics.Completed + scheduler.Metrics.Failed
		total := jobs

		if done == total {
			elapsed := time.Since(start).Seconds()

			result := BenchmarkResult{
				Workers:     workers,
				Jobs:        jobs,
				Seconds:     elapsed,
				Throughput:  float64(jobs) / elapsed,
				Utilization: scheduler.Metrics.Utilization,
				Completed:   scheduler.Metrics.Completed,
				Failed:      scheduler.Metrics.Failed,
			}

			benchmarkHistory = append(benchmarkHistory, result)
			scheduler.Mutex.Unlock()
			break
		}
		scheduler.Mutex.Unlock()
		time.Sleep(200 * time.Millisecond)
	}
}

var benchmarkHistory []BenchmarkResult
