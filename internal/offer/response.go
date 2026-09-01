package offer

type Response struct {
	ID           uint    `json:"id"`
	ListingID    uint    `json:"listing_id"`
	BuyerID      uint    `json:"buyer_id"`
	OfferedPrice float64 `json:"offered_price"`
	Quantity     float64 `json:"quantity"`
	Message      string  `json:"message"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

func ToResponse(o *Offer) Response {
	return Response{
		ID:           o.ID,
		ListingID:    o.ListingID,
		BuyerID:      o.BuyerID,
		OfferedPrice: o.OfferedPrice,
		Quantity:     o.Quantity,
		Message:      o.Message,
		Status:       o.Status,
		CreatedAt:    o.CreatedAt.Format("2006-01-02"),
	}
}

func ToResponseList(offers []Offer) []Response {
	res := make([]Response, 0, len(offers))
	for _, o := range offers {
		item := o
		res = append(res, ToResponse(&item))
	}
	return res
}