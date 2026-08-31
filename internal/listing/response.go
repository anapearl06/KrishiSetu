package listing

type Response struct {
	ID          uint    `json:"id"`
	FarmerID    uint    `json:"farmer_id"`
	Crop        string  `json:"crop"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	State       string  `json:"state"`
	District    string  `json:"district"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

func ToResponse(l *CropListing) Response {
	return Response{
		ID:          l.ID,
		FarmerID:    l.FarmerID,
		Crop:        l.CropName,
		Quantity:    l.Quantity,
		Unit:        l.Unit,
		Price:       l.ExpectedPrice,
		Description: l.Description,
		State:       l.State,
		District:    l.District,
		Status:      l.Status,
		CreatedAt:   l.CreatedAt.Format("2006-01-02"),
	}
}

func ToResponseList(listings []CropListing) []Response {
	res := make([]Response, 0, len(listings))
	for _, l := range listings {
		item := l
		res = append(res, ToResponse(&item))
	}
	return res
}
