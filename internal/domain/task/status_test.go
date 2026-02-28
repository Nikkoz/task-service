package task

import "testing"

func TestStatus_Valid(t *testing.T) {
	cases := []struct {
		name  string
		in    Status
		valid bool
	}{
		{"planned", StatusPlanned, true},
		{"in_progress", StatusInProgress, true},
		{"done", StatusDone, true},
		{"empty", Status(""), false},
		{"unknown", Status("lol"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.in.Valid(); got != tc.valid {
				t.Fatalf("expected %v, got %v", tc.valid, got)
			}
		})
	}
}

func TestNewStatus_IsValid(t *testing.T) {
	cases := []struct {
		name  string
		in    string
		valid bool
	}{
		{"planned", StatusPlanned.String(), true},
		{"in_progress", StatusInProgress.String(), true},
		{"done", StatusDone.String(), true},
		{"empty", "", false},
		{"unknown", "unknown", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewStatus(tc.in)
			if tc.valid {
				if err != nil || got.String() != tc.in {
					t.Fatalf("expected %v, got %v", tc.valid, got)
				}
			} else if err == nil {
				t.Fatalf("expected error, got %v", got)
			}
		})
	}
}
