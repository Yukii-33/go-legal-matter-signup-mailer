package main

import (
	"testing"
	"time"
)

func TestNeedsDeadlineFollowUp(t *testing.T) {
	now := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name     string
		deadline time.Time
		want     bool
	}{
		{"before deadline", now.Add(time.Hour), false},
		{"at deadline", now, true},
		{"past deadline", now.Add(-time.Hour), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NeedsDeadlineFollowUp(now, tc.deadline); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
