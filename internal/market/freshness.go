package market

import "time"

// calculateFreshness returns the age of the reported data
// and a human-readable freshness status.
//
// KrishiSetu intentionally does not claim that mandi data
// is real-time. The status is based on the reported arrival date.
func calculateFreshness(reportedDate time.Time) (int, string) {
	if reportedDate.IsZero() {
		return 0, "Unknown"
	}

	now := time.Now()

	// Compare dates rather than exact timestamps because
	// mandi arrival dates represent reporting dates.
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	reported := time.Date(
		reportedDate.Year(),
		reportedDate.Month(),
		reportedDate.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	age := int(today.Sub(reported).Hours() / 24)

	if age < 0 {
		age = 0
	}

	switch {
	case age == 0:
		return age, "Today"
	case age == 1:
		return age, "1 day old"
	case age <= 3:
		return age, "Recent"
	case age <= 7:
		return age, "1 week old"
	default:
		return age, "Older data"
	}
}