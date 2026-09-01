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
		Trend:        "Stable",
	}

	if len(records) == 1 {
		return result
	}

	previous := records[len(records)-2]

	result.PreviousPrice = previous.ModalPrice

	if previous.ModalPrice == 0 {
		return result
	}

	result.ChangePercent =
		((latest.ModalPrice - previous.ModalPrice) /
			previous.ModalPrice) * 100

	switch {
	case result.ChangePercent > 0:
		result.Trend = "Rising"

	case result.ChangePercent < 0:
		result.Trend = "Falling"

	default:
		result.Trend = "Stable"
	}

	return result
}