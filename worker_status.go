package main

type WorkerStatus struct {
	ID    int    `json:"id"`
	State string `json:"state"`
	JobID int    `json:"jobId"`
}
