package market

import (
	"testing"
	"time"
)

func TestCalculatePriceIntelligence_NilReturnsNil(t *testing.T) {
	result := CalculatePriceIntelligence(nil)
	if result != nil {
		t.Fatalf("expected nil for empty input, got %+v", result)
	}
}

func TestCalculatePriceIntelligence_SingleRecord(t *testing.T) {
	now := time.Now()
	records := []MarketPrice{
		{Commodity: "Wheat", ModalPrice: 2500, MinPrice: 2300, MaxPrice: 2700, ArrivalDate: now, Source: "mandi.in"},
	}
	result := CalculatePriceIntelligence(records)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Commodity != "Wheat" {
		t.Errorf("expected commodity Wheat, got %s", result.Commodity)
	}
	if result.CurrentPrice != 2500 {
		t.Errorf("expected price 2500, got %f", result.CurrentPrice)
	}
	if result.MinPrice != 2300 || result.MaxPrice != 2700 {
		t.Errorf("expected min 2300 max 2700, got %f/%f", result.MinPrice, result.MaxPrice)
	}
	if result.Source != "mandi.in" {
		t.Errorf("expected source mandi.in, got %s", result.Source)
	}
}

func TestCalculatePriceIntelligence_MultipleRecords(t *testing.T) {
	t1 := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	records := []MarketPrice{
		{Commodity: "Onion", ModalPrice: 1000, MinPrice: 800, MaxPrice: 1200, ArrivalDate: t1, Source: "mandiA"},
		{Commodity: "Onion", ModalPrice: 1200, MinPrice: 900, MaxPrice: 1500, ArrivalDate: t2, Source: "mandiB"},
	}
	result := CalculatePriceIntelligence(records)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// Average modal = (1000+1200)/2 = 1100
	if result.CurrentPrice != 1100 {
		t.Errorf("expected price 1100, got %f", result.CurrentPrice)
	}
	// Min should be 800, Max 1500
	if result.MinPrice != 800 || result.MaxPrice != 1500 {
		t.Errorf("expected min 800 max 1500, got %f/%f", result.MinPrice, result.MaxPrice)
	}
	// Latest date is t2, source mandiB
	if result.Source != "mandiB" {
		t.Errorf("expected source mandiB, got %s", result.Source)
	}
	if result.ReportedDate != "2026-09-03" {
		t.Errorf("expected reported_date 2026-09-03, got %s", result.ReportedDate)
	}
}

func TestCalculateFreshness_OldDate(t *testing.T) {
	oldDate := time.Now().AddDate(0, 0, -30)
	age, label := calculateFreshness(oldDate)
	if age < 29 {
		t.Errorf("expected age ~30, got %d", age)
	}
	if label != "Older data" {
		t.Errorf("expected Older data, got %s", label)
	}
}

func TestCalculateRegionalPriceComparison_Empty(t *testing.T) {
	result := CalculateRegionalPriceComparison(nil, "", "")
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}

func TestCalculateRegionalPriceComparison_GroupByState(t *testing.T) {
	now := time.Now()
	records := []MarketPrice{
		{Commodity: "Wheat", State: "UP", District: "Agra", Market: "A1", ModalPrice: 2000, MinPrice: 1800, MaxPrice: 2200, ArrivalDate: now, Source: "src"},
		{Commodity: "Wheat", State: "Punjab", District: "Ludhiana", Market: "P1", ModalPrice: 2100, MinPrice: 1900, MaxPrice: 2300, ArrivalDate: now, Source: "src2"},
	}
	result := CalculateRegionalPriceComparison(records, "", "")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestCalculateRegionalPriceComparison_GroupByDistrict(t *testing.T) {
	now := time.Now()
	records := []MarketPrice{
		{Commodity: "Wheat", State: "UP", District: "Agra", Market: "A1", ModalPrice: 2000, MinPrice: 1800, MaxPrice: 2200, ArrivalDate: now, Source: "src"},
		{Commodity: "Wheat", State: "UP", District: "Meerut", Market: "M1", ModalPrice: 2100, MinPrice: 1900, MaxPrice: 2300, ArrivalDate: now, Source: "src2"},
	}
	result := CalculateRegionalPriceComparison(records, "UP", "")
	if len(result) != 2 {
		t.Fatalf("expected 2 results (one per district), got %d", len(result))
	}
}

func TestCalculateRegionalPriceComparison_SpecificLocation(t *testing.T) {
	now := time.Now()
	records := []MarketPrice{
		{Commodity: "Wheat", State: "UP", District: "Agra", Market: "A1", ModalPrice: 2000, MinPrice: 1800, MaxPrice: 2200, ArrivalDate: now, Source: "src"},
		{Commodity: "Wheat", State: "UP", District: "Agra", Market: "A2", ModalPrice: 2100, MinPrice: 1900, MaxPrice: 2300, ArrivalDate: now, Source: "src"},
	}
	result := CalculateRegionalPriceComparison(records, "UP", "Agra")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].AverageModalPrice != 2050 {
		t.Errorf("expected avg 2050, got %f", result[0].AverageModalPrice)
	}
	if result[0].MarketCount != 2 {
		t.Errorf("expected 2 markets, got %d", result[0].MarketCount)
	}
}
