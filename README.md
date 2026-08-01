# ClearOutside Scraper

> An unofficial, Docker-ready Go service that retrieves publicly available Clear Outside forecasts, parses the HTML, and exposes the data through a simple JSON REST API.

This project is designed for personal and self-hosted use, providing a lightweight way to integrate Clear Outside forecast data into your own applications, dashboards, or home automation systems.

## Disclaimer

This project is **not affiliated with, endorsed by, or maintained by Clear Outside or First Light Optics Ltd.**

It retrieves information from **publicly accessible forecast pages** by parsing the HTML returned by the Clear Outside website. No private APIs or authenticated endpoints are used.

This software is intended for **personal, local, and self-hosted use only**. Users are responsible for ensuring their use complies with the upstream website's Terms of Use and applicable laws.

To minimize load on the upstream service, this project is designed to support response caching and should not be used to generate excessive automated traffic.

## License

MIT
