// Package scraper defines weather-provider contracts and Clear Outside parsing.
package scraper

import "context"

// Provider retrieves forecasts from a weather provider.
type Provider interface {
	Forecast(ctx context.Context, lat, lon float64) (*Forecast, error)
}

// Forecast is the provider-neutral REST forecast model.
type Forecast struct {
	Location  Location `json:"location"`
	Generated string   `json:"generated,omitempty"`
	Timezone  string   `json:"timezone,omitempty"`
	Days      []Day    `json:"days"`
}

// Location identifies a forecast location.
type Location struct {
	Name      string  `json:"name"`
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Day contains a local forecast date and hourly forecast values.
type Day struct {
	Date     string   `json:"date"`
	Moon     Moon     `json:"moon"`
	Sun      Sun      `json:"sun"`
	Darkness Darkness `json:"darkness"`
	Hours    []Hour   `json:"hours"`
}

// Moon contains lunar phase and rise/set information.
type Moon struct {
	Phase        string `json:"phase,omitempty"`
	Illumination string `json:"illumination,omitempty"`
	Rise         string `json:"rise,omitempty"`
	Set          string `json:"set,omitempty"`
}

// Sun contains solar rise, set, and transit times.
type Sun struct {
	Rise    string `json:"rise,omitempty"`
	Set     string `json:"set,omitempty"`
	Transit string `json:"transit,omitempty"`
}

// Darkness contains civil, nautical, and astronomical darkness intervals.
type Darkness struct {
	Civil        string `json:"civil,omitempty"`
	Nautical     string `json:"nautical,omitempty"`
	Astronomical string `json:"astronomical,omitempty"`
}

// Hour contains forecast conditions for one local hour.
type Hour struct {
	Hour          int           `json:"hour"`
	Rating        string        `json:"rating"`
	Cloud         Cloud         `json:"cloud"`
	Wind          Wind          `json:"wind"`
	Visibility    *float64      `json:"visibility_miles,omitempty"`
	Fog           *float64      `json:"fog_percent,omitempty"`
	Temperature   *float64      `json:"temperature_celsius,omitempty"`
	Humidity      *float64      `json:"humidity_percent,omitempty"`
	DewPoint      *float64      `json:"dew_point_celsius,omitempty"`
	Pressure      *float64      `json:"pressure_mb,omitempty"`
	Precipitation Precipitation `json:"precipitation"`
}

// Cloud contains sky obscuration percentages.
type Cloud struct {
	Total  *float64 `json:"total_percent,omitempty"`
	Low    *float64 `json:"low_percent,omitempty"`
	Medium *float64 `json:"medium_percent,omitempty"`
	High   *float64 `json:"high_percent,omitempty"`
}

// Wind contains speed and direction data.
type Wind struct {
	SpeedMPH  *float64 `json:"speed_mph,omitempty"`
	Direction string   `json:"direction,omitempty"`
	Degrees   *float64 `json:"degrees,omitempty"`
}

// Precipitation contains precipitation type, probability, and amount.
type Precipitation struct {
	Type        string   `json:"type,omitempty"`
	Probability *float64 `json:"probability_percent,omitempty"`
	AmountMM    *float64 `json:"amount_mm,omitempty"`
}
