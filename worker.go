package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, scheduler *Scheduler) {
	defer scheduler.WG.Done()

	fmt.Printf("Worker %d started\n", id)

	for {
		// job := <-scheduler.JobQueue
		select {
		case <-scheduler.Ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return

		case job := <-scheduler.JobQueue:

			scheduler.Mutex.Lock()
			scheduler.Workers[id].State = "Busy"
			scheduler.Workers[id].JobID = job.ID
			scheduler.Mutex.Unlock()

			fmt.Printf("Worker %d started Job %d\n", id, job.ID)

			// time.Sleep(job.Duration)
			ctx, cancel := context.WithTimeout(
				context.Background(),
				job.Timeout,
			)
			// defer cancel()

			err := ExecuteJob(ctx, job)

			cancel() // change made

			if err == context.DeadlineExceeded {
				fmt.Printf("⏰ Job %d timed out.\n", job.ID)

				scheduler.Mutex.Lock()
				scheduler.Workers[id].State = "Idle"
				scheduler.Workers[id].JobID = 0
				scheduler.Metrics.Failed++
				scheduler.Metrics.Queued--
				scheduler.Mutex.Unlock()

				continue
			}

			failed := rand.Intn(4) == 0
			if failed {

				job.Retries++
				if job.Retries <= job.MaxRetries {
					delay := time.Duration(1<<uint(job.Retries-1)) * time.Second

					fmt.Printf("❌ Worker %d: Job %d failed. Retrying in %v...\n", id, job.ID, delay)

					go func(j Job) {
						// time.Sleep(delay)
						// scheduler.Retry(j)
						select {
						case <-time.After(delay):
							scheduler.Retry(j)
						case <-scheduler.Ctx.Done():
							fmt.Printf("Retry for Job %d cancelled due to shutdown\n", j.ID)
							return
						}
					}(job)
				} else {
					fmt.Printf(
						"💥 Job %d permanently failed.\n", job.ID,
					)

					scheduler.Mutex.Lock()
					scheduler.Metrics.Failed++
					scheduler.Metrics.Queued--
					scheduler.Mutex.Unlock()
				}

				scheduler.Mutex.Lock()
				scheduler.Workers[id].State = "Idle"
				scheduler.Workers[id].JobID = 0
				scheduler.Mutex.Unlock()
				continue

			}
			fmt.Printf("✅ Worker %d: Job %d completed\n", id, job.ID)

			scheduler.Mutex.Lock()
			scheduler.Workers[id].State = "Idle"
			scheduler.Workers[id].JobID = 0
			scheduler.Metrics.Completed++
			scheduler.Metrics.Queued--
			scheduler.Mutex.Unlock()
		}

	}
}
