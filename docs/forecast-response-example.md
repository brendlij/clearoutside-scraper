# Forecast Response Example

Example response from `GET /forecast?lat=52.52&lon=13.40`. This is a representative, shortened response; production responses contain seven daily forecasts and hourly values for each day.

```json
{
  "location": {
    "name": "Mitte",
    "country": "Germany",
    "latitude": 52.52,
    "longitude": 13.4
  },
  "days": [
    {
      "date": "2026-08-01",
      "moon": {
        "phase": "Waning Gibbous",
        "illumination": "88%",
        "rise": "21:54",
        "set": "09:48"
      },
      "sun": {
        "rise": "05:28",
        "set": "20:58",
        "transit": "13:12"
      },
      "darkness": {
        "civil": "21:40 - 04:47",
        "nautical": "22:37 - 03:50",
        "astronomical": "00:02 - 02:28"
      },
      "hours": [
        {
          "hour": 12,
          "rating": "Bad",
          "cloud": {
            "total_percent": 100,
            "low_percent": 78,
            "medium_percent": 100,
            "high_percent": 87
          },
          "wind": {
            "speed_mph": 7,
            "direction": "North-West",
            "degrees": 311
          },
          "visibility_miles": 3,
          "fog_percent": 0,
          "temperature_celsius": 20,
          "humidity_percent": 91,
          "dew_point_celsius": 19,
          "pressure_mb": 1015,
          "precipitation": {
            "type": "Very Light Rain",
            "probability_percent": 78,
            "amount_mm": 0.3
          }
        },
        {
          "hour": 20,
          "rating": "OK",
          "cloud": {
            "total_percent": 38,
            "low_percent": 38,
            "medium_percent": 0,
            "high_percent": 0
          },
          "wind": {
            "speed_mph": 5,
            "direction": "North-East",
            "degrees": 36
          },
          "visibility_miles": 9,
          "fog_percent": 0,
          "temperature_celsius": 21,
          "humidity_percent": 70,
          "dew_point_celsius": 15,
          "pressure_mb": 1017,
          "precipitation": {
            "type": "None",
            "probability_percent": 12,
            "amount_mm": 0
          }
        },
        {
          "hour": 22,
          "rating": "Good",
          "cloud": {
            "total_percent": 14,
            "low_percent": 14,
            "medium_percent": 0,
            "high_percent": 0
          },
          "wind": {
            "speed_mph": 5,
            "direction": "East-North-East",
            "degrees": 67
          },
          "visibility_miles": 9,
          "fog_percent": 0,
          "temperature_celsius": 18,
          "humidity_percent": 78,
          "dew_point_celsius": 14,
          "pressure_mb": 1019,
          "precipitation": {
            "type": "None",
            "probability_percent": 3,
            "amount_mm": 0
          }
        }
      ]
    }
  ]
}
```