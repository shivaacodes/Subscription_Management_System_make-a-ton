package plans

import (
	"testing"
	"time"
)

func TestRecurrenceCalculation(t *testing.T) {
	service := NewService(nil)

	// Case 1: Daily
	jan1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	nextDaily := service.ComputeNextBillingDate(FrequencyDaily, jan1)
	if nextDaily.Day() != 2 || nextDaily.Month() != time.January {
		t.Errorf("expected Jan 2, got %v", nextDaily)
	}

	// Case 2: Weekly
	nextWeekly := service.ComputeNextBillingDate(FrequencyWeekly, jan1)
	if nextWeekly.Day() != 8 {
		t.Errorf("expected Jan 8, got %v", nextWeekly)
	}

	// Case 3: Monthly preserving month-end: Jan 31 -> Feb 28 (in non-leap 2026)
	jan31 := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	nextMonthly := service.ComputeNextBillingDate(FrequencyMonthly, jan31)
	if nextMonthly.Month() != time.February || nextMonthly.Day() != 28 {
		t.Errorf("expected Feb 28, got %v", nextMonthly)
	}

	// Case 4: Monthly leap year: Jan 31, 2024 -> Feb 29, 2024
	jan31_2024 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	nextLeap := service.ComputeNextBillingDate(FrequencyMonthly, jan31_2024)
	if nextLeap.Month() != time.February || nextLeap.Day() != 29 {
		t.Errorf("expected Feb 29, 2024, got %v", nextLeap)
	}

	// Case 5: Yearly preserving leap: Feb 29, 2024 -> Feb 28, 2025
	feb29_2024 := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	nextYearly := service.ComputeNextBillingDate(FrequencyYearly, feb29_2024)
	if nextYearly.Year() != 2025 || nextYearly.Month() != time.February || nextYearly.Day() != 28 {
		t.Errorf("expected Feb 28, 2025, got %v", nextYearly)
	}
}
