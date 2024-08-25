package repository

import (
	"time"
)

var (
	insertPriceQuery    = "INSERT INTO price (price, item_id, created_at) VALUES ($1, $2, $3)"
	selectAllItemsQuery = "SELECT id, name, url FROM item"
)

type PriceRecord struct {
	Price     float64
	ItemID    string
	CreatedAt time.Time
}

type TrackedItem struct {
	ID          string
	Name        string
	TrackingURL string
}

func GetAllItems() ([]*TrackedItem, error) {
	rows, err := connection.Query(selectAllItemsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*TrackedItem
	for rows.Next() {
		item := &TrackedItem{}
		err := rows.Scan(&item.ID, &item.Name, &item.TrackingURL)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func InsertPrice(itemID string, price float64) (*PriceRecord, error) {
	now := time.Now()
	_, err := connection.Exec(insertPriceQuery, price, itemID, now)
	if err != nil {
		return nil, err
	}
	return &PriceRecord{
		Price:     price,
		ItemID:    itemID,
		CreatedAt: now,
	}, nil
}
