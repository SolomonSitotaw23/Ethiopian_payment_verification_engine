package config

import "time"

type PerformanceConfig struct {
	MaxBatchSize       int
	DefaultConcurrency int
	Timeout            time.Duration
}

var Performance = PerformanceConfig{
	MaxBatchSize:       10,
	DefaultConcurrency: 10,
	Timeout:            60 * time.Second,
}
