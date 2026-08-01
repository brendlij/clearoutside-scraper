// Package config loads runtime settings from environment variables.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config contains application runtime settings.
type Config struct {
	Port          string
	CacheDuration time.Duration
	HTTPTimeout   time.Duration
	UserAgent     string
}

// Load returns configuration with safe defaults for local use.
func Load() Config {
	return Config{
		Port:          getenv("PORT", "2880"),
		CacheDuration: durationEnv("CACHE_DURATION", 30*time.Minute),
		HTTPTimeout:   durationEnv("HTTP_TIMEOUT", 15*time.Second),
		UserAgent:     getenv("USER_AGENT", "clearoutside-scraper/1.0 (personal self-hosted use)"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	if value, err := time.ParseDuration(os.Getenv(key)); err == nil && value > 0 {
		return value
	}
	if seconds, err := strconv.Atoi(os.Getenv(key)); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	return fallback
}
