package main

import (
	"log/slog"
	"net/http"
	"os"

	"clearoutside-scraper/internal/api"
	"clearoutside-scraper/internal/cache"
	"clearoutside-scraper/internal/config"
	"clearoutside-scraper/internal/scraper"
)

var version = "dev"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--healthcheck" {
		return
	}
	configuration := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	provider := scraper.NewClearOutsideProvider(&http.Client{Timeout: configuration.HTTPTimeout}, configuration.UserAgent, logger)
	handler := api.NewServer(provider, cache.NewMemory(configuration.CacheDuration), logger, version)
	server := &http.Server{Addr: ":" + configuration.Port, Handler: handler}
	logger.Info("server starting", "port", configuration.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
