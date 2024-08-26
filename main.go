package main

import (
	"price-tracker/logging"
	"price-tracker/scraper"
	"time"
)

func main() {
	duration := 6 * time.Hour
	logging.L.Printf("Starting scraper at fixed rate of %v", duration)
	scraper.RunAtFixedRate(duration)
}
