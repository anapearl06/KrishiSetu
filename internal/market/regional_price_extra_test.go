package market

import (
	"testing"
	"time"
)

// min/max across several markets are aggregated into a single region when the
// caller pins down a specific state + district.
func TestCalculateRegionalPriceComparison_MinMaxAcrossMarkets(t *testing.T) {
	now := time.Now()
	records := []MarketPrice{
		{State: "UP", District: "Agra", Market: "A1", ModalPrice: 2000, MinPrice: 1800, MaxPrice: 2200, ArrivalDate: now, Source: "src"},
		{State: "UP", District: "Agra", Market: "A2", ModalPrice: 2100, MinPrice: 1500, MaxPrice: 2400, ArrivalDate: now, Source: "src"},
		{State: "UP", District: "Agra", Market: "A3", ModalPrice: 2200, MinPrice: 1900, MaxPrice: 2300, ArrivalDate: now, Source: "src"},
	}
	result := CalculateRegionalPriceComparison(records, "UP", "Agra")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	// Minimum across all records = 1500, maximum = 2400.
	if result[0].MinPrice != 1500 {
		t.Errorf("expected min 1500, got %f", result[0].MinPrice)
	}
	if result[0].MaxPrice != 2400 {
		t.Errorf("expected max 2400, got %f", result[0].MaxPrice)
	}
	if result[0].MarketCount != 3 {
		t.Errorf("expected 3 markets, got %d", result[0].MarketCount)
	}
}

// The average modal price is computed over distinct markets, so two records
// from the same market count once toward the divisor.
func TestCalculateRegionalPriceComparison_AverageOverDistinctMarkets(t *testing.T) {
	now := time.Now()
	records := []MarketPrice{
		{State: "UP", District: "Agra", Market: "A1", ModalPrice: 2000, MinPrice: 1800, MaxPrice: 2200, ArrivalDate: now, Source: "src"},
		{State: "UP", District: "Agra", Market: "A1", ModalPrice: 2400, MinPrice: 1900, MaxPrice: 2500, ArrivalDate: now, Source: "src"},
		{State: "UP", District: "Agra", Market: "A2", ModalPrice: 3000, MinPrice: 2600, MaxPrice: 3200, ArrivalDate: now, Source: "src"},
	}
	result := CalculateRegionalPriceComparison(records, "UP", "Agra")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	// totalModal = 2000 + 2400 + 3000 = 7400, divided by 2 distinct markets = 3700.
	if result[0].AverageModalPrice != 3700 {
		t.Errorf("expected avg 3700, got %f", result[0].AverageModalPrice)
	}
	if result[0].MarketCount != 2 {
		t.Errorf("expected 2 distinct markets, got %d", result[0].MarketCount)
	}
}

// The newest arrival date's source is reported even when records arrive in
// different order within the loop.
func TestCalculateRegionalPriceComparison_UsesLatestSource(t *testing.T) {
	older := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	latest := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)
	records := []MarketPrice{
		{State: "UP", District: "Agra", Market: "A1", ModalPrice: 2000, MinPrice: 1800, MaxPrice: 2200, ArrivalDate: latest, Source: "mandiB"},
		{State: "UP", District: "Agra", Market: "A2", ModalPrice: 2100, MinPrice: 1900, MaxPrice: 2300, ArrivalDate: older, Source: "mandiA"},
	}
	result := CalculateRegionalPriceComparison(records, "UP", "Agra")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Source != "mandiB" {
		t.Errorf("expected source mandiB, got %s", result[0].Source)
	}
	if result[0].ReportedDate != "2026-09-05" {
		t.Errorf("expected reported date 2026-09-05, got %s", result[0].ReportedDate)
	}
}

// Freshness labels follow a progression from today through recent to older.
func TestCalculateFreshness_LabelProgression(t *testing.T) {
	today := time.Now()
	age, label := calculateFreshness(today)
	if age != 0 {
		t.Errorf("expected age 0 for today, got %d", age)
	}
	if label != "Today" {
		t.Errorf("expected Today, got %s", label)
	}

	recent := time.Now().AddDate(0, 0, -3)
	_, label = calculateFreshness(recent)
	if label != "Recent" {
		t.Errorf("expected Recent for 3-day-old data, got %s", label)
	}

	week := time.Now().AddDate(0, 0, -5)
	_, label = calculateFreshness(week)
	if label != "1 week old" {
		t.Errorf("expected 1 week old, got %s", label)
	}
}

// A zero arrival date yields an Unknown freshness status.
func TestCalculateFreshness_UnknownForZeroDate(t *testing.T) {
	age, label := calculateFreshness(time.Time{})
	if age != 0 {
		t.Errorf("expected age 0, got %d", age)
	}
	if label != "Unknown" {
		t.Errorf("expected Unknown, got %s", label)
	}
}
