# Njuškalo Crawler (Go + Playwright)

A specialized web scraper for Njuškalo.hr that finds VW Tiguans (or other cars) with specific filters:
- **Automatic Transmission** (DSG, Tiptronic, Sekvencijski)
- **Year**: 2014 - 2016
- **Price**: 10,000€ - 15,000€
- **Brand**: VW / Volkswagen

## Features
- **JavaScript Disabled**: Bypasses many behavioral bot detections.
- **Pagination**: Automatically crawls through all search results pages.
- **Persistence**: Remembers processed URLs in `processed_urls.json` to avoid duplicate work and redundant results.
- **Power Conversion**: Automatically converts engine power from **kW** to **HP**.
- **JSON Output**: Saves results to `results.json` and prints them to the console.

## Prerequisites
- [Go](https://golang.org/doc/install) (1.21 or higher recommended)
- [Playwright for Go](https://github.com/playwright-community/playwright-go)

## Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Install Chromium browser**:
   ```bash
   go run github.com/playwright-community/playwright-go/cmd/playwright install chromium
   ```

## Usage

### Build the binary
```bash
go build -o njuskalo-crawler .
```

### Run the crawler
```bash
./njuskalo-crawler
```

## How it works
The crawler visits the search page and collects all listing URLs. It then visits each individual ad page to extract detailed specifications like gearbox type, production year, and engine power. It applies strict filters in `filter.go` before presenting the final list.

## Disclaimer
This tool is for educational purposes only. Always respect the `robots.txt` and terms of service of the target website.
