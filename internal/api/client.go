package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"clearoutside-scraper/internal/cache"
	"clearoutside-scraper/internal/scraper"

	"github.com/go-chi/chi/v5"
)

// Server exposes provider-neutral forecast endpoints.
type Server struct {
	provider scraper.Provider
	cache    *cache.Memory
	logger   *slog.Logger
	version  string
}

// NewServer creates an HTTP handler for the supplied forecast provider.
func NewServer(provider scraper.Provider, memoryCache *cache.Memory, logger *slog.Logger, version string) http.Handler {
	server := &Server{provider: provider, cache: memoryCache, logger: logger, version: version}
	router := chi.NewRouter()
	router.Use(server.requestLogger)
	router.Get("/health", server.health)
	router.Get("/version", server.versionHandler)
	router.Get("/forecast", server.forecast)
	return router
}

func (server *Server) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, statusResponse{Status: "ok"})
}
func (server *Server) versionHandler(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, versionResponse{Version: server.version})
}

func (server *Server) forecast(writer http.ResponseWriter, request *http.Request) {
	lat, lon, err := coordinates(request)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "lat and lon must be valid coordinates"})
		return
	}
	key := fmt.Sprintf("%.6f:%.6f", lat, lon)
	if forecast, ok := server.cache.Get(key); ok {
		server.logger.Info("forecast cache hit", "key", key)
		writeJSON(writer, http.StatusOK, forecast)
		return
	}
	server.logger.Info("forecast cache miss", "key", key)
	started := time.Now()
	forecast, err := server.provider.Forecast(request.Context(), lat, lon)
	if err != nil {
		server.logger.Error("forecast retrieval failed", "error", err, "duration", time.Since(started))
		writeJSON(writer, http.StatusBadGateway, errorResponse{Error: "failed to retrieve forecast"})
		return
	}
	server.cache.Set(key, forecast)
	server.logger.Info("forecast scraped", "duration", time.Since(started))
	writeJSON(writer, http.StatusOK, forecast)
}

func coordinates(request *http.Request) (float64, float64, error) {
	lat, latErr := strconv.ParseFloat(request.URL.Query().Get("lat"), 64)
	lon, lonErr := strconv.ParseFloat(request.URL.Query().Get("lon"), 64)
	if latErr != nil || lonErr != nil || lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return 0, 0, fmt.Errorf("invalid coordinates")
	}
	return lat, lon, nil
}

func (server *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		server.logger.Info("incoming request", "method", request.Method, "path", request.URL.Path)
		next.ServeHTTP(writer, request)
	})
}
func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
