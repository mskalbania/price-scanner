package repository

import (
	"time"
)

var (
	insertPriceQuery    = "INSERT INTO price (price, item_id, created_at) VALUES ($1, $2, $3)"
	selectAllItemsQuery = `SELECT item.id, item.name, url, vendor.name, css_selector FROM item
							INNER JOIN source ON item.id = source.item_id
         					INNER JOIN vendor ON source.vendor_id = vendor.id`
)

type PriceRecord struct {
	Price     float64
	ItemID    string
	CreatedAt time.Time
}

type TrackedItem struct {
	ID          string
	Name        string
	URL         string
	Vendor      string
	CssSelector string
}

func GetAllItems() ([]TrackedItem, error) {
	rows, err := connection.Query(selectAllItemsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TrackedItem
	for rows.Next() {
		item := TrackedItem{}
		err := rows.Scan(&item.ID, &item.Name, &item.URL, &item.Vendor, &item.CssSelector)
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
