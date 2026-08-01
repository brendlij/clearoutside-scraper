package scraper

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const forecastBaseURL = "https://clearoutside.com/forecast"

// ClearOutsideProvider retrieves publicly available Clear Outside forecasts.
type ClearOutsideProvider struct {
	client    *http.Client
	userAgent string
	logger    *slog.Logger
}

// NewClearOutsideProvider creates a provider using client and userAgent.
func NewClearOutsideProvider(client *http.Client, userAgent string, logger *slog.Logger) *ClearOutsideProvider {
	return &ClearOutsideProvider{client: client, userAgent: userAgent, logger: logger}
}

// Forecast fetches and parses a forecast for the requested coordinates.
func (provider *ClearOutsideProvider) Forecast(ctx context.Context, lat, lon float64) (*Forecast, error) {
	fetchStarted := time.Now()
	body, err := provider.Fetch(ctx, lat, lon)
	if err != nil {
		provider.logger.Error("upstream forecast request failed", "error", err, "duration", time.Since(fetchStarted))
		return nil, err
	}
	defer body.Close()
	provider.logger.Info("forecast fetched", "duration", time.Since(fetchStarted))
	parseStarted := time.Now()
	forecast, err := Parse(body, lat, lon)
	if err != nil {
		provider.logger.Error("forecast parsing failed", "error", err, "duration", time.Since(parseStarted))
		return nil, err
	}
	provider.logger.Info("forecast parsed", "duration", time.Since(parseStarted))
	return forecast, nil
}

// Fetch downloads the raw forecast HTML.
func (provider *ClearOutsideProvider) Fetch(ctx context.Context, lat, lon float64) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", forecastBaseURL, url.PathEscape(fmt.Sprintf("%.6f", lat)), url.PathEscape(fmt.Sprintf("%.6f", lon)))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create forecast request: %w", err)
	}
	request.Header.Set("User-Agent", provider.userAgent)

	response, err := provider.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("fetch forecast: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("fetch forecast: upstream returned %s", response.Status)
	}
	return response.Body, nil
}
