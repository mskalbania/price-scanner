package main

import (
	"log"
	"price-tracker/scraper"
	"time"
)

func main() {
	duration := 3 * time.Hour
	log.Printf("Starting scraper at fixed rate of %v", duration)
	scraper.RunAtFixedRate(duration)
}
