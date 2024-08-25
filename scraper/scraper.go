package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"math/rand"
	"net/http"
	"price-tracker/repository"
	"strconv"
	"strings"
	"time"
)

var dbPath = "./db/price-db"
var client = &http.Client{
	Timeout: 5 * time.Second,
}

type scrapingResult struct {
	itemName  string
	itemPrice float64
}

func RunAtFixedRate(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for range ticker.C {
		repository.OpenConnection(dbPath)
		err := scrap()
		repository.CloseConnection()
		if err != nil {
			log.Println(err)
		}
	}
}

func scrap() error {
	items, err := repository.GetAllItems()
	if err != nil {
		return fmt.Errorf("error getting tracked items: %w", err)
	}

	results := make(chan *scrapingResult, len(items))
	group := &errgroup.Group{}

	for _, item := range items {
		group.Go(func() error {
			//random sleep between 50-250ms to avoid hitting the server at the same time
			time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)
			return scrapItemPrice(item, results)
		})
	}
	err = group.Wait()
	close(results)
	if err != nil {
		return err
	}
	if err := storeResults(results); err != nil {
		return err
	}
	return nil
}

func scrapItemPrice(item *repository.TrackedItem, results chan *scrapingResult) error {
	price, err := getPrice(item.TrackingURL)
	if err != nil {
		return fmt.Errorf("error getting price for %v: %w", item.Name, err)
	}
	results <- &scrapingResult{
		itemName:  item.Name,
		itemPrice: price,
	}
	return nil
}

func getPrice(url string) (float64, error) {
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
	return extractPrice(rs.Body)
}

func extractPrice(rs io.ReadCloser) (float64, error) {
	defer rs.Close()
	doc, err := goquery.NewDocumentFromReader(rs)
	if err != nil {
		return 0, fmt.Errorf("parsing html: %w", err)
	}
	priceText := doc.Find("#projector_price_value > span").Text()
	priceStr := strings.ReplaceAll(strings.Split(priceText, " ")[0], ",", ".")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing price: %w", err)
	}
	return price, nil
}

func storeResults(results chan *scrapingResult) error {
	fmt.Printf("\n---------------- Scrapping result [%v] ----------------\n", time.Now().Format("2006-01-02 15:04:05"))
	for result := range results {
		_, err := repository.InsertPrice(result.itemName, result.itemPrice)
		if err != nil {
			return fmt.Errorf("error inserting price record: %w", err)
		}
		fmt.Printf("%s - %.2f\n", result.itemName, result.itemPrice)
	}
	fmt.Println("------------------------------------------------------------------------")
	return nil
}
