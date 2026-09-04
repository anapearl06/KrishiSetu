package matching

type MatchResponse struct {
	ID             uint     `json:"id"`
	ListingID      uint     `json:"listing_id"`
	DemandID       uint     `json:"demand_id"`
	Score          float64  `json:"score"`
	Level          string   `json:"level"`
	CommodityScore float64  `json:"commodity_score"`
	QuantityScore  float64  `json:"quantity_score"`
	LocationScore  float64  `json:"location_score"`
	PriceScore     float64  `json:"price_score"`
	GradeScore     float64  `json:"grade_score"`
	Reasons        []string `json:"reasons"`
}

func ToResponse(match *Match) MatchResponse {
	return MatchResponse{
		ID:             match.ID,
		ListingID:      match.ListingID,
		DemandID:       match.DemandID,
		Score:          match.Score,
		Level:          match.Level,
		CommodityScore: match.CommodityScore,
		QuantityScore:  match.QuantityScore,
		LocationScore:  match.LocationScore,
		PriceScore:     match.PriceScore,
		GradeScore:     match.GradeScore,
		Reasons:        match.Reasons,
	}
}

func ToResponseList(matches []Match) []MatchResponse {
	responses := make([]MatchResponse, 0, len(matches))

	for i := range matches {
		responses = append(
			responses,
			ToResponse(&matches[i]),
		)
	}

	return responses
}