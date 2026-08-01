// Package cache provides short-lived forecast caching.
package cache

import (
	"sync"
	"time"

	"clearoutside-scraper/internal/scraper"
)

type entry struct {
	forecast *scraper.Forecast
	expires  time.Time
}

// Memory is a concurrency-safe, TTL-based forecast cache.
type Memory struct {
	mu      sync.RWMutex
	entries map[string]entry
	ttl     time.Duration
}

// NewMemory creates an in-memory cache with the given time-to-live.
func NewMemory(ttl time.Duration) *Memory {
	return &Memory{entries: make(map[string]entry), ttl: ttl}
}

// Get returns a non-expired forecast for key.
func (memory *Memory) Get(key string) (*scraper.Forecast, bool) {
	memory.mu.RLock()
	entry, ok := memory.entries[key]
	memory.mu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		if ok {
			memory.mu.Lock()
			delete(memory.entries, key)
			memory.mu.Unlock()
		}
		return nil, false
	}
	return entry.forecast, true
}

// Set stores forecast under key until its configured expiry.
func (memory *Memory) Set(key string, forecast *scraper.Forecast) {
	memory.mu.Lock()
	memory.entries[key] = entry{forecast: forecast, expires: time.Now().Add(memory.ttl)}
	memory.mu.Unlock()
}
