package query

import (
	"testing"
	"time"
)

func TestParse_Defaults(t *testing.T) {
	q, err := Parse(Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.TimeField != "time" {
		t.Errorf("expected default TimeField=\"time\", got %q", q.TimeField)
	}
	if q.Filter != nil {
		t.Error("expected nil Filter for empty expression")
	}
}

func TestParse_InvalidFilter(t *testing.T) {
	_, err := Parse(Options{FilterExpr: "!!!"})
	if err == nil {
		t.Fatal("expected error for invalid filter expression")
	}
}

func TestParse_InvalidSince(t *testing.T) {
	_, err := Parse(Options{Since: "not-a-date"})
	if err == nil {
		t.Fatal("expected error for invalid since value")
	}
}

func TestParse_ValidTimeRange(t *testing.T) {
	q, err := Parse(Options{
		Since: "2024-01-01T00:00:00Z",
		Until: "2024-12-31T23:59:59Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Since == nil || q.Until == nil {
		t.Fatal("expected non-nil Since and Until")
	}
}

func TestMatch_NoConstraints(t *testing.T) {
	q, _ := Parse(Options{})
	entry := Entry{"msg": "hello", "level": "info"}
	if !q.Match(entry) {
		t.Error("expected entry to match with no constraints")
	}
}

func TestMatch_TimeRange(t *testing.T) {
	now := time.Now().UTC()
	since := now.Add(-1 * time.Hour).Format(time.RFC3339)
	until := now.Add(1 * time.Hour).Format(time.RFC3339)

	q, err := Parse(Options{Since: since, Until: until})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inRange := Entry{"time": now.Format(time.RFC3339), "msg": "ok"}
	if !q.Match(inRange) {
		t.Error("expected in-range entry to match")
	}

	outOfRange := Entry{"time": now.Add(-2 * time.Hour).Format(time.RFC3339), "msg": "old"}
	if q.Match(outOfRange) {
		t.Error("expected out-of-range entry not to match")
	}
}

func TestMatch_MissingTimeField(t *testing.T) {
	q, _ := Parse(Options{Since: "2024-01-01T00:00:00Z"})
	entry := Entry{"msg": "no time field"}
	if q.Match(entry) {
		t.Error("expected entry without time field to be excluded when time bounds are set")
	}
}
