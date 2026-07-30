# Concurrent Job Scheduling System

A concurrent job scheduling system built in **Go** that efficiently executes background tasks using **goroutines**, **channels**, and a configurable **worker pool**.

The scheduler supports priority-based scheduling, retries with exponential backoff, timeout handling using `context.Context`, graceful shutdown, real-time worker monitoring, and benchmark-driven performance analysis.

A web dashboard provides live metrics, worker status, throughput, utilization, and benchmarking tools for comparing different worker pool configurations.

---

## Features

- Configurable worker pool
- Priority-based job scheduling using a heap
- Dispatcher-based scheduling architecture
- Retry mechanism with exponential backoff
- Job timeout and cancellation using `context.Context`
- Graceful shutdown without interrupting active jobs
- Live metrics dashboard
- Worker status monitoring
- Benchmarking with configurable worker and job counts
- Throughput and worker utilization tracking

---

## Features

* Configurable worker pool with dynamic concurrency levels
* Priority-based job scheduling using a heap-based priority queue
* Dispatcher goroutine for centralized job scheduling
* Buffered job queue for efficient worker distribution
* Concurrent job execution using goroutines
* Retry mechanism with exponential backoff
* Job timeout and cancellation using `context.Context`
* Graceful shutdown using `context.CancelFunc` and `sync.WaitGroup`
* Real-time worker status monitoring
* Live metrics dashboard displaying:

  * Completed jobs
  * Failed jobs
  * Queued jobs
  * Throughput
  * Worker utilization
* Configurable benchmarking with different worker and job counts
* Benchmark history table for comparing scheduler performance

## System Architecture

The scheduler follows a **Producer → Dispatcher → Worker Pool** architecture.

```text
                    +-------------------+
                    |   Job Submission  |
                    +---------+---------+
                              |
                              v
                  Priority Queue (Heap)
                              |
                    JobSignal Channel
                              |
                              v
                 Dispatcher Goroutine
                              |
                  Buffered Job Queue
                              |
          +---------+---------+---------+
          |         |         |         |
          v         v         v         v
      Worker 1  Worker 2  Worker 3 ... Worker N
          |         |         |         |
          +---------+---------+---------+
                              |
                              v
               Metrics & Worker Status
                              |
                              v
                  HTML / CSS / JS Dashboard
```

### Workflow

1. Jobs are submitted to the scheduler with configurable priority, timeout, retry count, and execution duration.
2. Submitted jobs are inserted into a priority queue implemented using Go's `container/heap`.
3. A dedicated dispatcher goroutine waits for job notifications and dispatches the highest-priority job into the buffered `JobQueue`.
4. Worker goroutines continuously listen on `JobQueue` and execute available jobs concurrently.
5. Each worker updates scheduler metrics and its current status after completing, failing, or retrying a job.
6. The frontend periodically fetches metrics and worker status through HTTP endpoints to provide a live dashboard.
7. Benchmark requests dynamically recreate the scheduler with the selected number of workers and jobs, allowing performance comparison across different configurations.

## Concurrency Model

The scheduler is designed around Go's concurrency primitives to efficiently process multiple independent jobs while maintaining thread safety and scalability.

### Goroutines

Concurrency is achieved using goroutines, which are lightweight threads managed by the Go runtime.

The scheduler creates:

* A dedicated **dispatcher goroutine** responsible for scheduling jobs.
* A configurable **worker pool**, where each worker runs as an independent goroutine and executes jobs concurrently.

This approach enables multiple jobs to be processed simultaneously while keeping the number of active execution threads under control.

---

### Channels

Channels are used for safe communication between goroutines.

#### JobSignal

The `JobSignal` channel is used to notify the dispatcher whenever a new job is submitted or retried. Rather than constantly polling the priority queue, the dispatcher remains blocked until it receives a signal, making scheduling efficient and event-driven.

#### JobQueue

The dispatcher removes the highest-priority job from the priority queue and places it into the buffered `JobQueue`. Workers continuously listen on this channel and receive jobs as soon as they become available.

Using channels eliminates the need for explicit synchronization while transferring jobs between concurrent components.

---

### Priority Queue

Jobs are stored in a heap-based priority queue implemented using Go's `container/heap` package.

Whenever multiple jobs are waiting, the dispatcher always selects the job with the highest priority before dispatching it to a worker.

This separates scheduling decisions from job execution while ensuring that more important tasks are processed first.

---

### Mutex Synchronization

Multiple goroutines access shared scheduler data simultaneously.

To prevent race conditions, a `sync.Mutex` protects shared resources including:

* Priority queue
* Scheduler metrics
* Worker status information

Only one goroutine can modify these resources at a time, ensuring data consistency throughout execution.

---

### Context-Based Cancellation

Each job executes with its own `context.Context` created using `context.WithTimeout()`.

If a job exceeds its configured timeout, the context is automatically cancelled and the worker marks the job as failed.

The scheduler also maintains a global context used during graceful shutdown to stop retries and notify running goroutines that shutdown has been initiated.

---

### WaitGroup

A `sync.WaitGroup` is used to coordinate graceful shutdown.

Before terminating the application, the scheduler waits until all worker goroutines complete their current execution, ensuring that no in-progress jobs are interrupted.

---

### Concurrent Monitoring

Worker status and scheduler metrics are continuously updated while jobs execute.

The frontend periodically retrieves this information through HTTP endpoints, providing real-time visibility into:

* Worker state
* Current job assignments
* Completed jobs
* Failed jobs
* Queue size
* Throughput
* Worker utilization

This allows the scheduler to be monitored without interrupting concurrent execution.

## Concurrency Model

The scheduler is designed around Go's concurrency primitives to efficiently process multiple independent jobs while maintaining thread safety and scalability.

### Goroutines

Concurrency is achieved using goroutines, which are lightweight threads managed by the Go runtime.

The scheduler creates:

* A dedicated **dispatcher goroutine** responsible for scheduling jobs.
* A configurable **worker pool**, where each worker runs as an independent goroutine and executes jobs concurrently.

This approach enables multiple jobs to be processed simultaneously while keeping the number of active execution threads under control.

---

### Channels

Channels are used for safe communication between goroutines.

#### JobSignal

The `JobSignal` channel is used to notify the dispatcher whenever a new job is submitted or retried. Rather than constantly polling the priority queue, the dispatcher remains blocked until it receives a signal, making scheduling efficient and event-driven.

#### JobQueue

The dispatcher removes the highest-priority job from the priority queue and places it into the buffered `JobQueue`. Workers continuously listen on this channel and receive jobs as soon as they become available.

Using channels eliminates the need for explicit synchronization while transferring jobs between concurrent components.

---

### Priority Queue

Jobs are stored in a heap-based priority queue implemented using Go's `container/heap` package.

Whenever multiple jobs are waiting, the dispatcher always selects the job with the highest priority before dispatching it to a worker.

This separates scheduling decisions from job execution while ensuring that more important tasks are processed first.

---

### Mutex Synchronization

Multiple goroutines access shared scheduler data simultaneously.

To prevent race conditions, a `sync.Mutex` protects shared resources including:

* Priority queue
* Scheduler metrics
* Worker status information

Only one goroutine can modify these resources at a time, ensuring data consistency throughout execution.

---

### Context-Based Cancellation

Each job executes with its own `context.Context` created using `context.WithTimeout()`.

If a job exceeds its configured timeout, the context is automatically cancelled and the worker marks the job as failed.

The scheduler also maintains a global context used during graceful shutdown to stop retries and notify running goroutines that shutdown has been initiated.

---

### WaitGroup

A `sync.WaitGroup` is used to coordinate graceful shutdown.

Before terminating the application, the scheduler waits until all worker goroutines complete their current execution, ensuring that no in-progress jobs are interrupted.

---

### Concurrent Monitoring

Worker status and scheduler metrics are continuously updated while jobs execute.

The frontend periodically retrieves this information through HTTP endpoints, providing real-time visibility into:

* Worker state
* Current job assignments
* Completed jobs
* Failed jobs
* Queue size
* Throughput
* Worker utilization

This allows the scheduler to be monitored without interrupting concurrent execution.

## Design Decisions

Several design choices were made to improve scalability, maintainability, and efficient resource utilization.

### Configurable Worker Pool

Instead of creating a new goroutine for every submitted job, the scheduler uses a configurable worker pool.

**Reasons:**

* Prevents uncontrolled goroutine creation.
* Limits resource consumption.
* Makes it easy to benchmark different concurrency levels by varying the number of workers.

---

### Dispatcher-Based Scheduling

A dedicated dispatcher goroutine is responsible for removing jobs from the priority queue and placing them into the buffered `JobQueue`.

**Reasons:**

* Centralizes scheduling logic.
* Prevents multiple workers from competing for the priority queue.
* Reduces synchronization complexity.
* Improves separation between scheduling and execution.

---

### Heap-Based Priority Queue

The scheduler stores jobs inside a heap-based priority queue using Go's `container/heap` package.

**Reasons:**

* Ensures higher-priority jobs are executed before lower-priority jobs.
* Provides efficient insertion and removal operations.
* Makes the scheduling policy independent of the worker implementation.

---

### Buffered Job Queue

Workers receive jobs through a buffered channel.

**Reasons:**

* Allows the dispatcher to continue dispatching jobs without immediately blocking.
* Improves throughput by reducing waiting time between scheduling and execution.
* Decouples scheduling from worker execution.

---

### Retry with Exponential Backoff

Failed jobs are retried after progressively increasing delays.

**Reasons:**

* Prevents repeated immediate retries.
* Reduces unnecessary resource usage during repeated failures.
* Simulates retry strategies commonly used in production systems.

---

### Context-Based Timeout

Every job executes with its own timeout using `context.WithTimeout()`.

**Reasons:**

* Prevents workers from being occupied indefinitely.
* Enables automatic cancellation of long-running jobs.
* Improves overall system responsiveness.

---

### Graceful Shutdown

The scheduler uses a shared cancellation context together with `sync.WaitGroup` to terminate cleanly.

**Reasons:**

* Running jobs are allowed to complete.
* Pending retries are cancelled.
* Worker goroutines exit safely.
* Prevents abrupt termination and resource leaks.

---

### Real-Time Monitoring

The frontend periodically retrieves scheduler metrics and worker status through HTTP endpoints.

**Reasons:**

* Provides live visibility into scheduler activity.
* Simplifies debugging and performance analysis.
* Allows benchmarking results to be observed in real time.

---

### Benchmark-Driven Evaluation

A benchmarking interface was included to compare scheduler performance under different workloads.

Users can vary:

* Number of workers
* Number of jobs

The scheduler records execution time, throughput, worker utilization, completed jobs, and failed jobs, making it possible to evaluate how concurrency impacts overall system performance.

## Benchmarking

The scheduler includes a built-in benchmarking interface to evaluate performance under different workloads.

Users can configure:

* **Worker Pool Size:** 1, 2, 4 or 8 workers
* **Number of Jobs:** 50, 100, 150 or 200 jobs

For each benchmark run, the dashboard records:

* Total execution time
* Throughput (jobs/second)
* Worker utilization
* Number of completed jobs
* Number of failed jobs

These results are displayed in a benchmark history table, making it easy to compare the impact of different worker pool sizes on overall scheduler performance.

> **Note:** Increasing the number of workers generally improves throughput until hardware resources or scheduling overhead become the limiting factors.

## Project Structure

```text
ConcurrentScheduler/
│
├── main.go                 # Application entry point
├── scheduler.go            # Scheduler implementation
├── worker.go               # Worker pool implementation
├── dispatcher.go           # Job dispatcher
├── priority_queue.go       # Heap-based priority queue
├── job.go                  # Job definition
├── job_executor.go         # Job execution logic
├── benchmark.go            # Benchmark execution
├── metrics.go              # Scheduler metrics
├── worker_status.go        # Worker status model
│
├── frontend/
│   ├── index.html
│   ├── style.css
│   └── script.js
│
└── README.md
```

## Running the Project

### Prerequisites

* Go 1.20 or later

### Clone the Repository

```bash
git clone https://github.com/dikshamisra18/ConcurrentScheduler.git
cd ConcurrentScheduler
```

### Run the Scheduler

```bash
go run .
```

Open your browser and navigate to:

```text
http://localhost:8080
```

From the dashboard, you can:

* Monitor live scheduler metrics
* View worker activity
* Configure benchmark parameters
* Compare benchmark results

## Screenshots

### Dashboard

*(Add dashboard screenshot here.)*

### Benchmark Results

*(Add benchmark results screenshot here.)*

