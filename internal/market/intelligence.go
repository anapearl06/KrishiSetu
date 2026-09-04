package market

type PriceIntelligence struct {
	Commodity string `json:"commodity"`

	CurrentPrice float64 `json:"current_price"`
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`

	PreviousPrice float64 `json:"previous_price"`
	ChangePercent float64 `json:"change_percent"`

	Trend string `json:"trend"`

	ReportedDate string `json:"reported_date"`
	Source       string `json:"source"`
}

func CalculatePriceIntelligence(
	records []MarketPrice,
) *PriceIntelligence {
	if len(records) == 0 {
		return nil
	}

	latest := records[len(records)-1]

	result := &PriceIntelligence{
		Commodity:    latest.Commodity,
		CurrentPrice: latest.ModalPrice,
		MinPrice:     latest.MinPrice,
		MaxPrice:     latest.MaxPrice,
		ReportedDate: latest.ArrivalDate.Format("2006-01-02"),
		Source:       latest.Source,
		Trend:        "Unavailable",
	}

	return result
}
