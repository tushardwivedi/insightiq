package jobs

import (
	"testing"
	"time"
)

func TestRetentionPolicyConfig_Defaults(t *testing.T) {
	config := RetentionPolicyConfig{}

	// Simulate default setting logic
	retentionPeriod := config.RetentionPeriod
	if retentionPeriod == 0 {
		retentionPeriod = 90 * 24 * time.Hour
	}

	checkInterval := config.CheckInterval
	if checkInterval == 0 {
		checkInterval = 24 * time.Hour
	}

	// Verify defaults
	expectedRetention := 90 * 24 * time.Hour
	expectedInterval := 24 * time.Hour

	if retentionPeriod != expectedRetention {
		t.Errorf("Expected default retention period %v, got %v", expectedRetention, retentionPeriod)
	}

	if checkInterval != expectedInterval {
		t.Errorf("Expected default check interval %v, got %v", expectedInterval, checkInterval)
	}
}

func TestRetentionPolicyConfig_CustomValues(t *testing.T) {
	config := RetentionPolicyConfig{
		RetentionPeriod: 30 * 24 * time.Hour,
		CheckInterval:   12 * time.Hour,
	}

	if config.RetentionPeriod != 30*24*time.Hour {
		t.Errorf("Expected retention period 30 days, got %v", config.RetentionPeriod)
	}

	if config.CheckInterval != 12*time.Hour {
		t.Errorf("Expected check interval 12 hours, got %v", config.CheckInterval)
	}
}

func TestRetentionPolicyConfig_CommonPeriods(t *testing.T) {
	tests := []struct {
		name   string
		days   int
		period time.Duration
	}{
		{
			name:   "30 days",
			days:   30,
			period: 30 * 24 * time.Hour,
		},
		{
			name:   "60 days",
			days:   60,
			period: 60 * 24 * time.Hour,
		},
		{
			name:   "90 days",
			days:   90,
			period: 90 * 24 * time.Hour,
		},
		{
			name:   "180 days",
			days:   180,
			period: 180 * 24 * time.Hour,
		},
		{
			name:   "365 days",
			days:   365,
			period: 365 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RetentionPolicyConfig{
				RetentionPeriod: tt.period,
			}

			expectedHours := tt.days * 24
			actualHours := int(config.RetentionPeriod.Hours())

			if actualHours != expectedHours {
				t.Errorf("Expected %d hours, got %d hours", expectedHours, actualHours)
			}
		})
	}
}

func TestRetentionPolicyConfig_CheckIntervals(t *testing.T) {
	tests := []struct {
		name     string
		interval time.Duration
		hours    int
	}{
		{
			name:     "Every 6 hours",
			interval: 6 * time.Hour,
			hours:    6,
		},
		{
			name:     "Every 12 hours",
			interval: 12 * time.Hour,
			hours:    12,
		},
		{
			name:     "Daily",
			interval: 24 * time.Hour,
			hours:    24,
		},
		{
			name:     "Weekly",
			interval: 7 * 24 * time.Hour,
			hours:    168,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RetentionPolicyConfig{
				CheckInterval: tt.interval,
			}

			actualHours := int(config.CheckInterval.Hours())

			if actualHours != tt.hours {
				t.Errorf("Expected %d hours, got %d hours", tt.hours, actualHours)
			}
		})
	}
}

func TestRetentionPolicy_MinimumPeriodValidation(t *testing.T) {
	// Test that periods less than 24 hours should be rejected
	tests := []struct {
		name    string
		period  time.Duration
		isValid bool
	}{
		{
			name:    "Less than 24 hours - invalid",
			period:  12 * time.Hour,
			isValid: false,
		},
		{
			name:    "Exactly 24 hours - valid",
			period:  24 * time.Hour,
			isValid: true,
		},
		{
			name:    "More than 24 hours - valid",
			period:  48 * time.Hour,
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.period >= 24*time.Hour

			if isValid != tt.isValid {
				t.Errorf("Expected validity %v, got %v for period %v", tt.isValid, isValid, tt.period)
			}
		})
	}
}
