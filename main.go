package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup
var scheduler *Scheduler

func metricsHandler(w http.ResponseWriter, r *http.Request) {

	scheduler.Mutex.Lock()
	metrics := scheduler.Metrics

	busy := 0
	for _, worker := range scheduler.Workers {
		if worker.State == "Busy" {
			busy++
		}
	}

	if scheduler.WorkerCount > 0 {
		metrics.Utilization = float64(busy) / float64(scheduler.WorkerCount)
	}

	elapsed := time.Since(scheduler.StartTime).Seconds()

	if elapsed > 0 {
		metrics.Throughput = float64(metrics.Completed) / elapsed
	}

	scheduler.Mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func workersHandler(w http.ResponseWriter, r *http.Request) {

	scheduler.Mutex.Lock()
	workers := make([]WorkerStatus, 0)

	for _, worker := range scheduler.Workers {
		workers = append(workers, *worker)
	}

	scheduler.Mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}

func benchmarkHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BenchmarkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// fmt.Printf("Benchmark requested: %d workers, %d jobs\n",
	// 	req.Workers,
	// 	req.Jobs,
	// )

	go StartBenchmark(req.Workers, req.Jobs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "started",
	})
}

func benchmarkResultsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(benchmarkHistory)
}

func main() {

	// jobs := make(chan Job)

	// go worker(1, jobs, &wg)
	// go worker(2, jobs, &wg)
	// go worker(3, jobs, &wg)

	scheduler = NewScheduler()

	go scheduler.Dispatch()

	// scheduler := NewScheduler()
	scheduler.StartWorkers(3)

	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		<-signals
		fmt.Println("\nGracefully Shutting Down...")
		scheduler.Cancel()
		scheduler.WG.Wait()
		os.Exit(0)
	}()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/workers", workersHandler)
	http.HandleFunc("/benchmark", benchmarkHandler)
	http.HandleFunc("/benchmarks", benchmarkResultsHandler)

	go func() {
		for i := 1; i <= 10; i++ {
			// wg.Add(1)
			job := Job{
				ID:         i,
				Priority:   rand.Intn(10) + 1,
				Duration:   time.Duration((i%3)+1) * time.Second,
				Timeout:    time.Duration(rand.Intn(4)+2) * time.Second,
				Retries:    0,
				MaxRetries: 3,
			}

			// if i == 1 {
			// 	job.Duration = 6 * time.Second
			// }
			// if i == 2 {
			// 	job.Priority = 8
			// }
			// if i == 6 {
			// 	job.Priority = 6
			// }
			// if i == 9 {
			// 	job.Priority = 10
			// }
			// if i == 5 {
			// 	job.Priority = 2
			// }

			fmt.Println("Submitting Job", i)
			// jobs <- job
			scheduler.Submit(job)
		}
	}()

	// wg.Wait()
	// fmt.Println("All jobs completed.")

	// time.Sleep(15 * time.Second)

	// time.Sleep(3 * time.Second)
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
