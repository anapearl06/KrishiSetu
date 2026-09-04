package matching

import (
	"strings"

	"github.com/raaj2493/KrishiSetu/internal/demand"
	"github.com/raaj2493/KrishiSetu/internal/listing"
)

const (
	commodityWeight = 30.0
	quantityWeight  = 20.0
	locationWeight  = 20.0
	priceWeight     = 15.0
	gradeWeight     = 15.0
)

type MatchResult struct {
	Score          float64  `json:"score"`
	Level          string   `json:"level"`
	CommodityScore float64  `json:"commodity_score"`
	QuantityScore  float64  `json:"quantity_score"`
	LocationScore  float64  `json:"location_score"`
	PriceScore     float64  `json:"price_score"`
	GradeScore     float64  `json:"grade_score"`
	Reasons        []string `json:"reasons"`
}

func CalculateMatch(
	cropListing listing.CropListing,
	buyerDemand demand.Demand,
) MatchResult {
	commodityScore, commodityReason :=
		calculateCommodityScore(
			cropListing.CropName,
			buyerDemand.CropName,
		)

	if commodityScore == 0 {
		return MatchResult{
			Score:          0,
			Level:          "Weak Match",
			CommodityScore: 0,
			QuantityScore:  0,
			LocationScore:  0,
			PriceScore:     0,
			GradeScore:     0,
			Reasons: []string{
				"Commodity does not match",
			},
		}
	}

	quantityScore, quantityReason :=
		calculateQuantityScore(
			cropListing.Quantity,
			cropListing.Unit,
			buyerDemand.Quantity,
			buyerDemand.Unit,
		)

	locationScore, locationReason :=
		calculateLocationScore(
			cropListing.State,
			cropListing.District,
			buyerDemand.State,
			buyerDemand.District,
		)

	priceScore, priceReason :=
		calculatePriceScore(
			cropListing.ExpectedPrice,
			buyerDemand.TargetPrice,
		)

	gradeScore, gradeReason :=
		calculateGradeScore(cropListing.QualityGrade)

	total :=
		commodityScore +
			quantityScore +
			locationScore +
			priceScore +
			gradeScore

	if total > 100 {
		total = 100
	}

	if total < 0 {
		total = 0
	}

	reasons := []string{
		commodityReason,
		quantityReason,
		locationReason,
		priceReason,
		gradeReason,
	}

	return MatchResult{
		Score:          total,
		Level:          calculateMatchLevel(total),
		CommodityScore: commodityScore,
		QuantityScore:  quantityScore,
		LocationScore:  locationScore,
		PriceScore:     priceScore,
		GradeScore:     gradeScore,
		Reasons:        reasons,
	}
}

func calculateCommodityScore(
	listingCrop string,
	demandCrop string,
) (float64, string) {
	listingCrop = normalize(listingCrop)
	demandCrop = normalize(demandCrop)

	if listingCrop == "" || demandCrop == "" {
		return 0, "Commodity information is missing"
	}

	if listingCrop == demandCrop {
		return commodityWeight, "Same commodity"
	}

	return 0, "Commodity does not match"
}

func calculateQuantityScore(
	listingQuantity float64,
	listingUnit string,
	demandQuantity float64,
	demandUnit string,
) (float64, string) {
	if listingQuantity <= 0 || demandQuantity <= 0 {
		return 0, "Quantity information is invalid"
	}

	if normalize(listingUnit) != normalize(demandUnit) {
		return 0, "Quantity units are different"
	}

	if listingQuantity >= demandQuantity {
		return quantityWeight, "Required quantity available"
	}

	ratio := listingQuantity / demandQuantity

	if ratio <= 0 {
		return 0, "Required quantity is not available"
	}

	score := ratio * quantityWeight

	if score > quantityWeight {
		score = quantityWeight
	}

	return score, "Part of the required quantity is available"
}

func calculateLocationScore(
	listingState string,
	listingDistrict string,
	demandState string,
	demandDistrict string,
) (float64, string) {
	listingState = normalize(listingState)
	listingDistrict = normalize(listingDistrict)

	demandState = normalize(demandState)
	demandDistrict = normalize(demandDistrict)

	if listingState == "" || demandState == "" {
		return 0, "Location information is missing"
	}

	if listingState == demandState &&
		listingDistrict != "" &&
		listingDistrict == demandDistrict {
		return locationWeight, "Same district"
	}

	if listingState == demandState {
		return 15, "Same state"
	}

	return 5, "Different state"
}

func calculatePriceScore(
	expectedPrice float64,
	targetPrice float64,
) (float64, string) {
	if expectedPrice <= 0 || targetPrice <= 0 {
		return 0, "Price information is invalid"
	}

	if expectedPrice == targetPrice {
		return priceWeight, "Price exactly matches buyer target"
	}

	difference := abs(expectedPrice - targetPrice)
	percentageDifference := difference / targetPrice

	switch {
	case percentageDifference <= 0.05:
		return 13, "Price is within 5% of buyer target"

	case percentageDifference <= 0.10:
		return 10, "Price is within 10% of buyer target"

	case percentageDifference <= 0.20:
		return 6, "Price is within 20% of buyer target"

	default:
		return 0, "Price is far from buyer target"
	}
}

func calculateGradeScore(
	listingGrade string,
) (float64, string) {
	listingGrade = normalize(listingGrade)

	if listingGrade == "" {
		return gradeWeight,
			"Buyer has no specific grade requirement"
	}

	return gradeWeight,
		"Quality grade information available"
}

func calculateMatchLevel(score float64) string {
	switch {
	case score >= 90:
		return "Excellent Match"

	case score >= 75:
		return "Good Match"

	case score >= 60:
		return "Possible Match"

	default:
		return "Weak Match"
	}
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}

	return value
}
