# ClearOutside Scraper

> An unofficial, Docker-ready Go service that retrieves publicly available Clear Outside forecasts, parses the HTML, and exposes the data through a simple JSON REST API.

This project is designed for personal and self-hosted use, providing a lightweight way to integrate Clear Outside forecast data into your own applications, dashboards, or home automation systems.

## Disclaimer

This project is **not affiliated with, endorsed by, or maintained by Clear Outside or First Light Optics Ltd.**

It retrieves information from **publicly accessible forecast pages** by parsing the HTML returned by the Clear Outside website. No private APIs or authenticated endpoints are used.

This software is intended for **personal, local, and self-hosted use only**. Users are responsible for ensuring their use complies with the upstream website's Terms of Use and applicable laws.

To minimize load on the upstream service, this project is designed to support response caching and should not be used to generate excessive automated traffic.

## Run

Requires Go 1.25 or later.

```powershell
go run ./cmd/server
```

Or run the Docker image:

```powershell
docker compose up --build
```

The service listens on port `2880` by default.

## API

```text
GET /health
GET /version
GET /forecast?lat=52.52&lon=13.40
```

`/forecast` returns a provider-neutral JSON document containing location details, seven daily forecasts, hourly ratings, cloud coverage, wind, visibility, fog, temperature, humidity, dew point, pressure, precipitation, moon, sun, and darkness details.

See the [forecast response example](docs/forecast-response-example.md) for a representative JSON payload.

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `2880` | HTTP listener port. |
| `CACHE_DURATION` | `30m` | Forecast cache lifetime. Go duration strings or seconds are accepted. |
| `HTTP_TIMEOUT` | `15s` | Upstream Clear Outside request timeout. Go duration strings or seconds are accepted. |
| `USER_AGENT` | `clearoutside-scraper/1.0 (personal self-hosted use)` | User-Agent sent to Clear Outside. |

## License

MIT
