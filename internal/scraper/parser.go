package scraper

import (
	"fmt"
	"html"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var locationPattern = regexp.MustCompile(`^Forecast for (.+?) \(([-+0-9.]+),\s*([-+0-9.]+)\)$`)
var numberPattern = regexp.MustCompile(`[-+]?\d+(?:\.\d+)?`)
var windPattern = regexp.MustCompile(`([0-9.]+)mph from the (.+?) \(([0-9.]+).+?\)`)
var rangePattern = regexp.MustCompile(`Forecast:\s*(\d{2}/\d{2}/\d{2})\s+to\s+(\d{2}/\d{2}/\d{2})`)

// Parse transforms Clear Outside forecast HTML into a provider-neutral model.
func Parse(source io.Reader, lat, lon float64) (*Forecast, error) {
	document, err := goquery.NewDocumentFromReader(source)
	if err != nil {
		return nil, fmt.Errorf("parse forecast HTML: %w", err)
	}
	forecast := &Forecast{Location: Location{Latitude: lat, Longitude: lon}}
	if err := parseLocation(document, &forecast.Location); err != nil {
		return nil, err
	}
	dates, err := forecastDates(document)
	if err != nil {
		return nil, err
	}
	var parseErr error
	dayIndex := 0
	document.Find("#forecast .fc_day").Each(func(_ int, selection *goquery.Selection) {
		if parseErr != nil {
			return
		}
		day, dayErr := parseDay(selection)
		if dayErr != nil {
			parseErr = dayErr
			return
		}
		if dayIndex >= len(dates) {
			parseErr = fmt.Errorf("parse forecast: found more days than forecast date range")
			return
		}
		day.Date = dates[dayIndex]
		dayIndex++
		forecast.Days = append(forecast.Days, day)
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if len(forecast.Days) == 0 {
		return nil, fmt.Errorf("parse forecast: no .fc_day elements found")
	}
	return forecast, nil
}

func forecastDates(document *goquery.Document) ([]string, error) {
	matches := rangePattern.FindStringSubmatch(normalizedText(document.Find("h2").First()))
	if len(matches) != 3 {
		return nil, fmt.Errorf("parse forecast dates: expected forecast date range in heading")
	}
	start, err := time.Parse("02/01/06", matches[1])
	if err != nil {
		return nil, fmt.Errorf("parse forecast start date: %w", err)
	}
	end, err := time.Parse("02/01/06", matches[2])
	if err != nil {
		return nil, fmt.Errorf("parse forecast end date: %w", err)
	}
	if end.Before(start) {
		return nil, fmt.Errorf("parse forecast dates: end date precedes start date")
	}
	var dates []string
	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date.Format("2006-01-02"))
	}
	return dates, nil
}

func parseLocation(document *goquery.Document, location *Location) error {
	text := normalizedText(document.Find("h1").First())
	matches := locationPattern.FindStringSubmatch(text)
	if len(matches) != 4 {
		return fmt.Errorf("parse location: expected forecast heading, got %q", text)
	}
	coordinatesLat, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return fmt.Errorf("parse location latitude: %w", err)
	}
	coordinatesLon, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return fmt.Errorf("parse location longitude: %w", err)
	}
	parts := strings.SplitN(matches[1], ",", 2)
	location.Name = strings.TrimSpace(parts[0])
	if len(parts) == 2 {
		location.Country = strings.TrimSpace(parts[1])
	}
	location.Latitude, location.Longitude = coordinatesLat, coordinatesLon
	return nil
}

func parseDay(selection *goquery.Selection) (Day, error) {
	date := normalizedText(selection.Find(".fc_day_date").First())
	if date == "" {
		return Day{}, fmt.Errorf("parse day: missing .fc_day_date")
	}
	ratings := selection.Find(".fc_hour_ratings li")
	if ratings.Length() == 0 {
		return Day{}, fmt.Errorf("parse day %q: no hourly ratings", date)
	}
	day := Day{Hours: make([]Hour, ratings.Length())}
	ratings.Each(func(index int, rating *goquery.Selection) {
		day.Hours[index] = Hour{Hour: integerText(rating), Rating: normalizedText(rating.Find("span").Last())}
	})
	rows := detailRows(selection)
	parseClouds(rows, day.Hours)
	parseWind(rows, day.Hours)
	parseMoon(selection, &day.Moon)
	parseSunAndDarkness(selection, &day.Sun, &day.Darkness)
	parseNumericRow(rows, "Visibility (miles)", day.Hours, func(hour *Hour, value *float64) { hour.Visibility = value })
	parseNumericRow(rows, "Fog (%)", day.Hours, func(hour *Hour, value *float64) { hour.Fog = value })
	parseNumericRow(rows, "Temperature (°C)", day.Hours, func(hour *Hour, value *float64) { hour.Temperature = value })
	parseNumericRow(rows, "Relative Humidity (%)", day.Hours, func(hour *Hour, value *float64) { hour.Humidity = value })
	parseNumericRow(rows, "Dew Point (°C)", day.Hours, func(hour *Hour, value *float64) { hour.DewPoint = value })
	parseNumericRow(rows, "Pressure (mb)", day.Hours, func(hour *Hour, value *float64) { hour.Pressure = value })
	parsePrecipitation(rows, day.Hours)
	return day, nil
}

func parseMoon(day *goquery.Selection, moon *Moon) {
	moon.Phase = normalizedText(day.Find(".fc_moon_phase").First())
	moon.Illumination = normalizedText(day.Find(".fc_moon_percentage").First())
	times := strings.Fields(normalizedText(day.Find(".fc_moon_riseset").First()))
	if len(times) >= 2 {
		moon.Rise, moon.Set = times[0], times[1]
	}
}

func parseSunAndDarkness(day *goquery.Selection, sun *Sun, darkness *Darkness) {
	content := day.Find(".fc_daylight").First().AttrOr("data-content", "")
	sun.Rise = popupValue(content, "Sunrise")
	sun.Set = popupValue(content, "Sunset")
	sun.Transit = popupValue(content, "Sun Transit")
	darkness.Civil = popupValue(content, "Civil Dark")
	darkness.Nautical = popupValue(content, "Nautical Dark")
	darkness.Astronomical = popupValue(content, "Astro Dark")
}

func popupValue(content, label string) string {
	pattern := regexp.MustCompile(`(?i)<strong>` + regexp.QuoteMeta(label) + `:</strong>\s*([^<]+)`)
	matches := pattern.FindStringSubmatch(html.UnescapeString(content))
	if len(matches) != 2 {
		return ""
	}
	return strings.TrimSpace(strings.ReplaceAll(matches[1], "\u00a0", " "))
}

func detailRows(day *goquery.Selection) map[string]*goquery.Selection {
	rows := make(map[string]*goquery.Selection)
	day.Find(".fc_detail_row").Each(func(_ int, row *goquery.Selection) {
		label := normalizedText(row.Find(".fc_detail_label").First())
		if label != "" {
			rows[label] = row
		}
	})
	return rows
}

func parseClouds(rows map[string]*goquery.Selection, hours []Hour) {
	parseNumericRow(rows, "Total Clouds (% Sky Obscured)", hours, func(hour *Hour, value *float64) { hour.Cloud.Total = value })
	parseNumericRow(rows, "Low Clouds (% Sky Obscured)", hours, func(hour *Hour, value *float64) { hour.Cloud.Low = value })
	parseNumericRow(rows, "Medium Clouds (% Sky Obscured)", hours, func(hour *Hour, value *float64) { hour.Cloud.Medium = value })
	parseNumericRow(rows, "High Clouds (% Sky Obscured)", hours, func(hour *Hour, value *float64) { hour.Cloud.High = value })
}

func parseWind(rows map[string]*goquery.Selection, hours []Hour) {
	forEachCell(rows["Wind Speed/Direction (mph)"], hours, func(hour *Hour, cell *goquery.Selection) {
		matches := windPattern.FindStringSubmatch(cell.AttrOr("title", ""))
		if len(matches) != 4 {
			return
		}
		hour.Wind.SpeedMPH = floatPointer(matches[1])
		hour.Wind.Direction = matches[2]
		hour.Wind.Degrees = floatPointer(matches[3])
	})
}

func parsePrecipitation(rows map[string]*goquery.Selection, hours []Hour) {
	forEachCell(rows["Precipitation Type"], hours, func(hour *Hour, cell *goquery.Selection) {
		hour.Precipitation.Type = cell.AttrOr("title", normalizedText(cell))
	})
	parseNumericRow(rows, "Precipitation Probability (%)", hours, func(hour *Hour, value *float64) { hour.Precipitation.Probability = value })
	parseNumericRow(rows, "Precipitation Amount (mm)", hours, func(hour *Hour, value *float64) { hour.Precipitation.AmountMM = value })
}

func parseNumericRow(rows map[string]*goquery.Selection, label string, hours []Hour, assign func(*Hour, *float64)) {
	forEachCell(rows[label], hours, func(hour *Hour, cell *goquery.Selection) { assign(hour, floatPointer(normalizedText(cell))) })
}

func forEachCell(row *goquery.Selection, hours []Hour, apply func(*Hour, *goquery.Selection)) {
	if row == nil {
		return
	}
	row.Find(".fc_hours li").EachWithBreak(func(index int, cell *goquery.Selection) bool {
		if index >= len(hours) {
			return false
		}
		apply(&hours[index], cell)
		return true
	})
}

func integerText(selection *goquery.Selection) int {
	value, _ := strconv.Atoi(numberPattern.FindString(normalizedText(selection)))
	return value
}
func floatPointer(value string) *float64 {
	parsed, err := strconv.ParseFloat(numberPattern.FindString(value), 64)
	if err != nil {
		return nil
	}
	return &parsed
}
func normalizedText(selection *goquery.Selection) string {
	return strings.Join(strings.Fields(selection.Text()), " ")
}
