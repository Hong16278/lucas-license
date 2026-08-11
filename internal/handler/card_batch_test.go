package handler

import "testing"

func TestLicenseKeyPrefix(t *testing.T) {
	tests := map[string]string{
		"automate":               "AUTOMATE",
		"auto-mate_pro":          "AUTOMATEPRO",
		"软件":                     "KG",
		"very-long-product-slug": "VERYLONGPROD",
	}
	for input, want := range tests {
		if got := licenseKeyPrefix(input); got != want {
			t.Fatalf("licenseKeyPrefix(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestPlanDurationBounds(t *testing.T) {
	if err := validatePlanNumericBounds(planBounds{DurationDays: 36500}); err != nil {
		t.Fatalf("expected maximum duration to be accepted: %v", err)
	}
	if err := validatePlanNumericBounds(planBounds{DurationDays: -1}); err == nil {
		t.Fatal("expected negative duration to be rejected")
	}
	if err := validatePlanNumericBounds(planBounds{DurationDays: 36501}); err == nil {
		t.Fatal("expected excessive duration to be rejected")
	}
}
