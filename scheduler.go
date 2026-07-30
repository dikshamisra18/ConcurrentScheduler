package main

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	// JobQueue    chan Job
	Jobs        PriorityQueue
	JobQueue    chan Job
	JobSignal   chan struct{}
	WorkerCount int
	Metrics     Metrics
	Workers     map[int]*WorkerStatus
	Mutex       sync.Mutex
	WG          sync.WaitGroup
	Ctx         context.Context
	Cancel      context.CancelFunc
	StartTime   time.Time
}

func NewScheduler() *Scheduler {

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		// JobQueue:    make(chan Job),
		Jobs:        make(PriorityQueue, 0),
		JobQueue:    make(chan Job),
		JobSignal:   make(chan struct{}, 1),
		WorkerCount: 0,
		Metrics:     Metrics{},
		Workers:     make(map[int]*WorkerStatus),
		Ctx:         ctx,
		Cancel:      cancel,
		StartTime:   time.Now(),
	}
}

func (s *Scheduler) Submit(job Job) {

	s.Mutex.Lock()
	heap.Push(&s.Jobs, job)
	s.Metrics.Queued++
	s.Mutex.Unlock()

	// s.JobQueue <- job
	select {
	case s.JobSignal <- struct{}{}:
	default:
	}
}

func (s *Scheduler) StartWorkers(count int) {

	s.WorkerCount = count

	for i := 1; i <= count; i++ {
		s.WG.Add(1)
		s.Workers[i] = &WorkerStatus{
			ID:    i,
			State: "Idle",
			JobID: 0,
		}
		go worker(i, s)
	}
}

func (s *Scheduler) Retry(job Job) {
	// s.JobQueue <- job

	s.Mutex.Lock()
	heap.Push(&s.Jobs, job)
	s.Mutex.Unlock()

	select {
	case s.JobSignal <- struct{}{}:
	default:
	}
}

func (s *Scheduler) Dispatch() {
	for {
		select {
		case <-s.Ctx.Done():
			return
		case <-s.JobSignal:
			for {

				s.Mutex.Lock()

				if len(s.Jobs) == 0 {
					s.Mutex.Unlock()
					break
				}

				job := heap.Pop(&s.Jobs).(Job)
				s.Mutex.Unlock()

				s.JobQueue <- job
			}
		}
	}
}
