package market

import "time"

// RegionalPriceComparison represents aggregated price information
// for a specific region.
type RegionalPriceComparison struct {
	Region             string    `json:"region"`
	State              string    `json:"state"`
	District           string    `json:"district,omitempty"`
	AverageModalPrice  float64   `json:"average_modal_price"`
	MinPrice           float64   `json:"min_price"`
	MaxPrice           float64   `json:"max_price"`
	MarketCount        int       `json:"market_count"`
	ReportedDate       string    `json:"reported_date"`
	Source             string    `json:"source"`
	DataAgeDays        int       `json:"data_age_days"`
	Freshness          string    `json:"freshness"`
}

// CalculateRegionalPriceComparison aggregates mandi records.
//
// If a state is provided and district is empty, the records are
// grouped by district within that state.
//
// If neither state nor district is provided, the records are
// grouped by state.
//
// If both state and district are provided, a single regional
// aggregate is returned.
func CalculateRegionalPriceComparison(
	records []MarketPrice,
	state string,
	district string,
) []RegionalPriceComparison {

	if len(records) == 0 {
		return []RegionalPriceComparison{}
	}

	type aggregate struct {
		state       string
		district    string
		totalModal  float64
		minPrice    float64
		maxPrice    float64
		markets     map[string]struct{}
		reported    time.Time
		source      string
	}

	aggregates := make(map[string]*aggregate)

	for _, record := range records {
		var key string

		switch {
		case state == "" && district == "":
			// Compare states.
			key = record.State

		case state != "" && district == "":
			// Compare districts inside the selected state.
			key = record.District

		default:
			// Specific state + district.
			key = record.State + "|" + record.District
		}

		if _, exists := aggregates[key]; !exists {
			aggregates[key] = &aggregate{
				state:      record.State,
				district:   record.District,
				minPrice:   record.MinPrice,
				maxPrice:   record.MaxPrice,
				markets:    make(map[string]struct{}),
				reported:   record.ArrivalDate,
				source:     record.Source,
			}
		}

		agg := aggregates[key]

		agg.totalModal += record.ModalPrice

		if record.MinPrice < agg.minPrice {
			agg.minPrice = record.MinPrice
		}

		if record.MaxPrice > agg.maxPrice {
			agg.maxPrice = record.MaxPrice
		}

		agg.markets[record.Market] = struct{}{}

		if record.ArrivalDate.After(agg.reported) {
			agg.reported = record.ArrivalDate
			agg.source = record.Source
		}
	}

	results := make([]RegionalPriceComparison, 0, len(aggregates))

	for _, agg := range aggregates {
		var region string

		switch {
		case state == "" && district == "":
			region = agg.state

		case state != "" && district == "":
			region = agg.district

		default:
			region = agg.district
		}

		averageModal := agg.totalModal / float64(len(agg.markets))

		dataAgeDays, freshness := calculateFreshness(agg.reported)

		results = append(results, RegionalPriceComparison{
			Region:            region,
			State:             agg.state,
			District:          agg.district,
			AverageModalPrice: averageModal,
			MinPrice:          agg.minPrice,
			MaxPrice:          agg.maxPrice,
			MarketCount:       len(agg.markets),
			ReportedDate:      agg.reported.Format("2006-01-02"),
			Source:            agg.source,
			DataAgeDays:       dataAgeDays,
			Freshness:         freshness,
		})
	}

	return results
}