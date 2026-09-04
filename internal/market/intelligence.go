package market

type PriceIntelligence struct {
	Commodity     string  `json:"commodity"`
	CurrentPrice  float64 `json:"current_price"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	PreviousPrice float64 `json:"previous_price"`
	ChangePercent float64 `json:"change_percent"`
	Trend         string  `json:"trend"`
	ReportedDate  string  `json:"reported_date"`
	Source        string  `json:"source"`
	DataAgeDays   int     `json:"data_age_days"`
	Freshness     string  `json:"freshness"`
}

func CalculatePriceIntelligence(
	records []MarketPrice,
) *PriceIntelligence {

	if len(records) == 0 {
		return nil
	}

	var (
		totalModal float64
		minPrice   = records[0].MinPrice
		maxPrice   = records[0].MaxPrice
	)

	latestDate := records[0].ArrivalDate
	source := records[0].Source

	for _, record := range records {
		totalModal += record.ModalPrice

		if record.MinPrice < minPrice {
			minPrice = record.MinPrice
		}

		if record.MaxPrice > maxPrice {
			maxPrice = record.MaxPrice
		}

		if record.ArrivalDate.After(latestDate) {
			latestDate = record.ArrivalDate
			source = record.Source
		}
	}

	averageModal := totalModal / float64(len(records))

	dataAgeDays, freshness := calculateFreshness(latestDate)

	return &PriceIntelligence{
		Commodity:     records[0].Commodity,
		CurrentPrice:  averageModal,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		PreviousPrice: 0,
		ChangePercent: 0,
		Trend:         "Unavailable",
		ReportedDate:  latestDate.Format("2006-01-02"),
		Source:        source,
		DataAgeDays:   dataAgeDays,
		Freshness:     freshness,
	}
}