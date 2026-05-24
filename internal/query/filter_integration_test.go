package query

import (
	"testing"
	"time"
)

func strPtr(s string) *string { return &s }

func TestCompile_NoFilter(t *testing.T) {
	q := Query{}
	cq, err := Compile(q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cq.Matcher != nil {
		t.Error("expected nil matcher when no filter provided")
	}
}

func TestCompile_ValidFilter(t *testing.T) {
	q := Query{Filter: strPtr(`level == "error"`)}
	cq, err := Compile(q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cq.Matcher == nil {
		t.Fatal("expected non-nil matcher")
	}
}

func TestCompile_InvalidFilter(t *testing.T) {
	q := Query{Filter: strPtr(`level ===`)}
	_, err := Compile(q)
	if err == nil {
		t.Fatal("expected error for invalid filter")
	}
}

func TestCompiledQuery_Match_FilterAndTime(t *testing.T) {
	now := time.Now().UTC()
	since := now.Add(-1 * time.Hour)
	until := now.Add(1 * time.Hour)

	q := Query{
		Since: &since,
		Until: &until,
		Filter: strPtr(`level == "error"`)}

	cq, err := Compile(q)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	entry := map[string]any{
		"time":  now.Format(time.RFC3339),
		"level": "error",
		"msg":   "something failed",
	}
	if !cq.Match(entry) {
		t.Error("expected entry to match")
	}

	entry["level"] = "info"
	if cq.Match(entry) {
		t.Error("expected info entry to not match error filter")
	}
}

func TestCompiledQuery_Match_OutOfTimeRange(t *testing.T) {
	now := time.Now().UTC()
	since := now.Add(1 * time.Hour) // future

	q := Query{Since: &since}
	cq, err := Compile(q)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	entry := map[string]any{
		"time":  now.Format(time.RFC3339),
		"level": "error",
	}
	if cq.Match(entry) {
		t.Error("expected entry outside time range to not match")
	}
}
