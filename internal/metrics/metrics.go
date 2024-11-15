package metrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/projectname/internal/logger"
)

type Metrics struct {
	counters map[string]int64
	timers   map[string]time.Duration
	logger   *logger.Logger
	mu       sync.RWMutex
}

var (
	instance *Metrics
	once     sync.Once
)

func New() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			counters: make(map[string]int64),
			timers:   make(map[string]time.Duration),
			logger:   logger.New(),
		}
	})
	return instance
}

func (m *Metrics) IncrementCounter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name]++
}

func (m *Metrics) GetCounter(name string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counters[name]
}

func (m *Metrics) RecordTime(name string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timers[name] = duration
}

func (m *Metrics) GetTimer(name string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.timers[name]
}

// Example middleware for tracking request metrics
func (m *Metrics) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.IncrementCounter("http_requests_total")

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		m.RecordTime("http_request_duration", duration)
	})
}
