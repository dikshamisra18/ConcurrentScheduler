package main

type Metrics struct {
	Completed   int     `json:"completed"`
	Failed      int     `json:"failed"`
	Queued      int     `json:"queued"`
	Throughput  float64 `json:"throughput"`
	Utilization float64 `json:"utilization"`
}
