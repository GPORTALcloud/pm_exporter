package metric

import (
	"sync"
	"time"
)

type Storage struct {
	Mutex           *sync.Mutex
	persistDuration time.Duration
	Metrics         map[string]Metric
}

type Metric struct {
	Name string
	Type string
	Help string
	Vars map[string]MetricVar
}

type MetricVar struct {
	Value     interface{}
	ExpiresAt int64
}
