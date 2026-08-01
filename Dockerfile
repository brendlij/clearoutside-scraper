FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /clearoutside-scraper ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /clearoutside-scraper /clearoutside-scraper
EXPOSE 2880
USER nonroot:nonroot
ENTRYPOINT ["/clearoutside-scraper"]
