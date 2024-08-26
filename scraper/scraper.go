package scraper

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"math/rand"
	"net/http"
	"price-tracker/logging"
	"price-tracker/repository"
	"strconv"
	"strings"
	"sync"
	"time"
)

var dbPath = "./db/price-db"
var client = &http.Client{
	Timeout: 5 * time.Second,
}

type scrapingResult struct {
	itemName  string
	itemPrice float64
	vendor    string
}

func RunAtFixedRate(duration time.Duration) {
	for {
		repository.OpenConnection(dbPath)
		err := scrap()
		repository.CloseConnection()
		if err != nil {
			logging.L.Println(err)
		}
		time.Sleep(duration)
	}
}

func scrap() error {
	items, err := repository.GetAllItems()
	if err != nil {
		return fmt.Errorf("error getting tracked items: %w", err)
	}

	resultChan := make(chan scrapingResult, len(items))
	errorsChan := make(chan error, len(items))

	var wg sync.WaitGroup
	start := time.Now()

	for _, item := range items {
		wg.Add(1)
		go func() {
			//random sleep between 50-250ms to avoid hitting the server at the same time
			time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)
			scrapItemPrice(item, resultChan, errorsChan)
			wg.Done()
		}()
	}
	wg.Wait()
	timeTaken := time.Since(start)

	close(resultChan)
	close(errorsChan)
	errorArr := make([]error, 0)

	if err := storeResults(timeTaken, resultChan); err != nil {
		errorArr = append(errorArr, err)
	}
	for err := range errorsChan {
		errorArr = append(errorArr, err)
	}
	if len(errorArr) > 0 {
		return errors.Join(errorArr...)
	}
	return nil
}

func scrapItemPrice(item repository.TrackedItem, results chan scrapingResult, errors chan error) {
	price, err := getPrice(item.URL, item.CssSelector)
	if err != nil {
		errors <- fmt.Errorf("error getting price for %v: %w", item.Name, err)
	}
	results <- scrapingResult{
		itemName:  item.Name,
		itemPrice: price,
		vendor:    item.Vendor,
	}
}

func getPrice(url string, selector string) (float64, error) {
	rq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request for %v: %w", url, err)
	}
	rq.Header.Set("User-Agent", "Mozilla/5.0")
	rs, err := client.Do(rq)
	if err != nil {
		return 0, fmt.Errorf("sending request at %v: %w", url, err)
	}
	if rs.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("response status code from %v: %v", url, rs.StatusCode)
	}
	return extractPrice(selector, rs.Body)
}

func extractPrice(selector string, rs io.ReadCloser) (float64, error) {
	defer rs.Close()
	doc, err := goquery.NewDocumentFromReader(rs)
	if err != nil {
		return 0, fmt.Errorf("parsing html: %w", err)
	}
	priceText := doc.Find(selector).Text()
	priceStr := strings.ReplaceAll(strings.Split(priceText, " ")[0], ",", ".")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing price: %w", err)
	}
	return price, nil
}

func storeResults(timeTaken time.Duration, results chan scrapingResult) error {
	//todo add item grouping by name
	fmt.Printf("\n%s / %s\n", time.Now().Format("2006-01-02 15:04:05"), timeTaken.Round(time.Millisecond))
	for result := range results {
		_, err := repository.InsertPrice(result.itemName, result.itemPrice)
		if err != nil {
			return fmt.Errorf("error inserting price record: %w", err)
		}
		fmt.Printf("[%s] %s - %.2f\n", result.vendor, result.itemName, result.itemPrice)
	}
	fmt.Println()
	return nil
}
