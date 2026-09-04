package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/raaj2493/KrishiSetu/internal/market"
)

type MandiRecord struct {
	State       string  `json:"state"`
	District    string  `json:"district"`
	Market      string  `json:"market"`
	Commodity   string  `json:"commodity"`
	Variety     string  `json:"variety"`
	Grade       string  `json:"grade"`
	ArrivalDate string  `json:"arrival_date"`
	MinPrice    float64 `json:"min_price"`
	MaxPrice    float64 `json:"max_price"`
	ModalPrice  float64 `json:"modal_price"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		fmt.Println("DATABASE_URL is not set")
		os.Exit(1)
	}

	file, err := os.Open("dataingestion/data/raw/mandi_prices.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var records []MandiRecord

	if err := json.NewDecoder(file).Decode(&records); err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d records from JSON\n", len(records))

	db, err := gorm.Open(
		postgres.Open(databaseURL),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	const batchSize = 500

	for start := 0; start < len(records); start += batchSize {
		end := start + batchSize

		if end > len(records) {
			end = len(records)
		}

		batch := records[start:end]

		prices := make([]market.MarketPrice, 0, len(batch))

		for _, record := range batch {
			arrivalDate, err := time.Parse(
				"02/01/2006",
				strings.TrimSpace(record.ArrivalDate),
			)
			if err != nil {
				panic(fmt.Errorf(
					"invalid arrival date %q: %w",
					record.ArrivalDate,
					err,
				))
			}

			prices = append(prices, market.MarketPrice{
				State:       strings.TrimSpace(record.State),
				District:    strings.TrimSpace(record.District),
				Market:      strings.TrimSpace(record.Market),
				Commodity:   strings.TrimSpace(record.Commodity),
				Variety:     strings.TrimSpace(record.Variety),
				Grade:       strings.TrimSpace(record.Grade),
				ArrivalDate: arrivalDate,
				MinPrice:    record.MinPrice,
				MaxPrice:    record.MaxPrice,
				ModalPrice:  record.ModalPrice,
				Source:      "data.gov.in",
			})
		}

		if err := db.CreateInBatches(&prices, batchSize).Error; err != nil {
			panic(fmt.Errorf(
				"failed inserting batch %d-%d: %w",
				start,
				end,
				err,
			))
		}

		fmt.Printf(
			"Inserted %d/%d records\n",
			end,
			len(records),
		)
	}

	fmt.Println("Mandi data import completed successfully.")
}