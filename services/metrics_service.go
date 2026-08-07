package services

import (
	"sync/atomic"
	"time"
)

type MetricsService struct {
	startTime      time.Time
	totalRequests  uint64
	validReceipts  uint64
	failedReceipts uint64
}

var Metrics = &MetricsService{
	startTime: time.Now(),
}

func (m *MetricsService) RecordVerification(success bool) {
	atomic.AddUint64(&m.totalRequests, 1)
	if success {
		atomic.AddUint64(&m.validReceipts, 1)
	} else {
		atomic.AddUint64(&m.failedReceipts, 1)
	}
}

type MetricsSnapshot struct {
	UptimeSeconds float64 `json:"uptimeSeconds"`
	TotalRequests uint64  `json:"totalRequests"`
	ValidReceipts uint64  `json:"validReceipts"`
	FailedReceipts uint64 `json:"failedReceipts"`
}

func (m *MetricsService) GetSnapshot() MetricsSnapshot {
	return MetricsSnapshot{
		UptimeSeconds:  time.Since(m.startTime).Seconds(),
		TotalRequests:  atomic.LoadUint64(&m.totalRequests),
		ValidReceipts:  atomic.LoadUint64(&m.validReceipts),
		FailedReceipts: atomic.LoadUint64(&m.failedReceipts),
	}
}
