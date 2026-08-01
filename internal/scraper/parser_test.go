package scraper

import (
	"os"
	"testing"
)

func TestParseSampleForecast(t *testing.T) {
	file, err := os.Open("../../sample.html")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	forecast, err := Parse(file, 52.52, 13.40)
	if err != nil {
		t.Fatal(err)
	}
	if forecast.Location.Name != "Mitte" || forecast.Location.Country != "Germany" {
		t.Fatalf("unexpected location: %+v", forecast.Location)
	}
	if len(forecast.Days) != 7 || len(forecast.Days[0].Hours) != 24 {
		t.Fatalf("unexpected forecast shape: %d days, %d hours", len(forecast.Days), len(forecast.Days[0].Hours))
	}
	first := forecast.Days[0].Hours[0]
	if first.Rating != "Bad" || first.Cloud.Total == nil || *first.Cloud.Total != 100 {
		t.Fatalf("expected first rating and cloud data, got %+v", first)
	}
}
