package matching

import ("testing" 

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing")

func TestCalculateCommodityScore(t *testing.T) {
	tests := []struct {
		name          string
		listingCrop   string
		demandCrop    string
		expectedScore float64
	}{
		{
			name:          "exact match",
			listingCrop:   "Wheat",
			demandCrop:    "Wheat",
			expectedScore: 30,
		},
		{
			name:          "case insensitive match",
			listingCrop:   "WHEAT",
			demandCrop:    "wheat",
			expectedScore: 30,
		},
		{
			name:          "whitespace is ignored",
			listingCrop:   " Wheat ",
			demandCrop:    "wheat",
			expectedScore: 30,
		},
		{
			name:          "different commodity",
			listingCrop:   "Wheat",
			demandCrop:    "Rice",
			expectedScore: 0,
		},
		{
			name:          "missing listing commodity",
			listingCrop:   "",
			demandCrop:    "Wheat",
			expectedScore: 0,
		},
		{
			name:          "missing demand commodity",
			listingCrop:   "Wheat",
			demandCrop:    "",
			expectedScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculateCommodityScore(
				tt.listingCrop,
				tt.demandCrop,
			)

			if score != tt.expectedScore {
				t.Fatalf(
					"expected score %.2f, got %.2f",
					tt.expectedScore,
					score,
				)
			}
		})
	}
}


func TestCalculateQuantityScore(t *testing.T) {
	tests := []struct {
		name          string
		listingQty    float64
		listingUnit   string
		demandQty     float64
		demandUnit    string
		expectedScore float64
	}{
		{
			name:          "listing fully satisfies demand",
			listingQty:    600,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 20,
		},
		{
			name:          "listing exactly satisfies demand",
			listingQty:    500,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 20,
		},
		{
			name:          "listing partially satisfies demand",
			listingQty:    250,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 10,
		},
		{
			name:          "listing provides small portion",
			listingQty:    100,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 4,
		},
		{
			name:          "zero listing quantity",
			listingQty:    0,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 0,
		},
		{
			name:          "negative listing quantity",
			listingQty:    -100,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 0,
		},
		{
			name:          "zero demand quantity",
			listingQty:    500,
			listingUnit:   "kg",
			demandQty:     0,
			demandUnit:    "kg",
			expectedScore: 0,
		},
		{
			name:          "different units",
			listingQty:    500,
			listingUnit:   "kg",
			demandQty:     500,
			demandUnit:    "ton",
			expectedScore: 0,
		},
		{
			name:          "unit comparison ignores case",
			listingQty:    500,
			listingUnit:   "KG",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 20,
		},
		{
			name:          "unit comparison ignores whitespace",
			listingQty:    500,
			listingUnit:   " kg ",
			demandQty:     500,
			demandUnit:    "kg",
			expectedScore: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculateQuantityScore(
				tt.listingQty,
				tt.listingUnit,
				tt.demandQty,
				tt.demandUnit,
			)

			if score != tt.expectedScore {
				t.Fatalf(
					"expected score %.2f, got %.2f",
					tt.expectedScore,
					score,
				)
			}
		})
	}
}

func TestCalculateLocationScore(t *testing.T) {
	tests := []struct {
		name          string
		listingState  string
		listingDistrict string
		demandState   string
		demandDistrict string
		expectedScore float64
	}{
		{
			name:            "same district",
			listingState:    "Uttar Pradesh",
			listingDistrict: "Ghaziabad",
			demandState:     "Uttar Pradesh",
			demandDistrict:  "Ghaziabad",
			expectedScore:   20,
		},
		{
			name:            "same state different district",
			listingState:    "Uttar Pradesh",
			listingDistrict: "Ghaziabad",
			demandState:     "Uttar Pradesh",
			demandDistrict:  "Meerut",
			expectedScore:   15,
		},
		{
			name:            "different state",
			listingState:    "Uttar Pradesh",
			listingDistrict: "Ghaziabad",
			demandState:     "Punjab",
			demandDistrict:  "Ludhiana",
			expectedScore:   5,
		},
		{
			name:            "missing listing state",
			listingState:    "",
			listingDistrict: "Ghaziabad",
			demandState:     "Uttar Pradesh",
			demandDistrict:  "Ghaziabad",
			expectedScore:   0,
		},
		{
			name:            "missing demand state",
			listingState:    "Uttar Pradesh",
			listingDistrict: "Ghaziabad",
			demandState:     "",
			demandDistrict:  "Ghaziabad",
			expectedScore:   0,
		},
		{
			name:            "case insensitive",
			listingState:    "UTTAR PRADESH",
			listingDistrict: "GHAZIABAD",
			demandState:     "uttar pradesh",
			demandDistrict:  "ghaziabad",
			expectedScore:   20,
		},
		{
			name:            "whitespace ignored",
			listingState:    " Uttar Pradesh ",
			listingDistrict: " Ghaziabad ",
			demandState:     "Uttar Pradesh",
			demandDistrict:  "Ghaziabad",
			expectedScore:   20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculateLocationScore(
				tt.listingState,
				tt.listingDistrict,
				tt.demandState,
				tt.demandDistrict,
			)

			if score != tt.expectedScore {
				t.Fatalf(
					"expected score %.2f, got %.2f",
					tt.expectedScore,
					score,
				)
			}
		})
	}
}


func TestCalculatePriceScore(t *testing.T) {
	tests := []struct {
		name          string
		expectedPrice float64
		targetPrice   float64
		expectedScore float64
	}{
		{
			name:          "exact price match",
			expectedPrice: 25,
			targetPrice:   25,
			expectedScore: 15,
		},
		{
			name:          "within five percent",
			expectedPrice: 26,
			targetPrice:   25,
			expectedScore: 13,
		},
		{
			name:          "within ten percent",
			expectedPrice: 27,
			targetPrice:   25,
			expectedScore: 10,
		},
		{
			name:          "within twenty percent",
			expectedPrice: 29,
			targetPrice:   25,
			expectedScore: 6,
		},
		{
			name:          "more than twenty percent",
			expectedPrice: 32,
			targetPrice:   25,
			expectedScore: 0,
		},
		{
			name:          "lower price within five percent",
			expectedPrice: 24,
			targetPrice:   25,
			expectedScore: 13,
		},
		{
			name:          "zero expected price",
			expectedPrice: 0,
			targetPrice:   25,
			expectedScore: 0,
		},
		{
			name:          "zero target price",
			expectedPrice: 25,
			targetPrice:   0,
			expectedScore: 0,
		},
		{
			name:          "negative expected price",
			expectedPrice: -10,
			targetPrice:   25,
			expectedScore: 0,
		},
		{
			name:          "negative target price",
			expectedPrice: 25,
			targetPrice:   -10,
			expectedScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculatePriceScore(
				tt.expectedPrice,
				tt.targetPrice,
			)

			if score != tt.expectedScore {
				t.Fatalf(
					"expected score %.2f, got %.2f",
					tt.expectedScore,
					score,
				)
			}
		})
	}
}

func TestCalculateGradeScore(t *testing.T) {
	tests := []struct {
		name          string
		listingGrade  string
		expectedScore float64
	}{
		{
			name:          "grade A available",
			listingGrade:  "A",
			expectedScore: 15,
		},
		{
			name:          "grade B available",
			listingGrade:  "B",
			expectedScore: 15,
		},
		{
			name:          "grade C available",
			listingGrade:  "C",
			expectedScore: 15,
		},
		{
			name:          "grade with whitespace",
			listingGrade:  " A ",
			expectedScore: 15,
		},
		{
			name:          "missing grade",
			listingGrade:  "",
			expectedScore: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculateGradeScore(tt.listingGrade)

			if score != tt.expectedScore {
				t.Fatalf(
					"expected score %.2f, got %.2f",
					tt.expectedScore,
					score,
				)
			}
		})
	}
}

func TestCalculateMatch(t *testing.T) {
	listing := listing.CropListing{
		CropName:       "Wheat",
		Quantity:       600,
		Unit:           "kg",
		ExpectedPrice:  26,
		QualityGrade:   "A",
		State:          "Uttar Pradesh",
		District:       "Ghaziabad",
	}

	demand := demand.Demand{
		CropName:     "Wheat",
		Quantity:     500,
		Unit:         "kg",
		TargetPrice:  25,
		State:        "Uttar Pradesh",
		District:     "Ghaziabad",
	}

	result := CalculateMatch(listing, demand)

	if result.Score != 98 {
		t.Fatalf("expected score 98, got %.2f", result.Score)
	}

	if result.Level != "Excellent Match" {
		t.Fatalf(
			"expected level %q, got %q",
			"Excellent Match",
			result.Level,
		)
	}

	if result.CommodityScore != 30 {
		t.Fatalf("expected commodity score 30, got %.2f", result.CommodityScore)
	}

	if result.QuantityScore != 20 {
		t.Fatalf("expected quantity score 20, got %.2f", result.QuantityScore)
	}

	if result.LocationScore != 20 {
		t.Fatalf("expected location score 20, got %.2f", result.LocationScore)
	}

	if result.PriceScore != 13 {
		t.Fatalf("expected price score 13, got %.2f", result.PriceScore)
	}

	if result.GradeScore != 15 {
		t.Fatalf("expected grade score 15, got %.2f", result.GradeScore)
	}

	if len(result.Reasons) != 5 {
		t.Fatalf(
			"expected 5 reasons, got %d",
			len(result.Reasons),
		)
	}
}

func TestCalculateMatchCommodityMismatch(t *testing.T) {
	listing := listing.CropListing{
		CropName:       "Wheat",
		Quantity:       500,
		Unit:           "kg",
		ExpectedPrice:  25,
		QualityGrade:   "A",
		State:          "Uttar Pradesh",
		District:       "Ghaziabad",
	}

	demand := demand.Demand{
		CropName:     "Rice",
		Quantity:     500,
		Unit:         "kg",
		TargetPrice:  25,
		State:        "Uttar Pradesh",
		District:     "Ghaziabad",
	}

	result := CalculateMatch(listing, demand)

	if result.Score != 0 {
		t.Fatalf("expected score 0, got %.2f", result.Score)
	}

	if result.Level != "Weak Match" {
		t.Fatalf(
			"expected level %q, got %q",
			"Weak Match",
			result.Level,
		)
	}

	if len(result.Reasons) != 1 {
		t.Fatalf(
			"expected 1 reason, got %d",
			len(result.Reasons),
		)
	}

	if result.Reasons[0] != "Commodity does not match" {
		t.Fatalf(
			"unexpected reason: %s",
			result.Reasons[0],
		)
	}
}