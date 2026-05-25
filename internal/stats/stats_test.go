package stats_test

import (
	"testing"
	"time"

	"logslice/internal/stats"
)

func TestCounter_InitialSummaryIsZero(t *testing.T) {
	c := stats.New()
	s := c.Summary()
	if s.Total != 0 || s.Matched != 0 || s.Skipped != 0 || s.Errored != 0 {
		t.Fatalf("expected zero counts, got %+v", s)
	}
}

func TestCounter_IncTotal(t *testing.T) {
	c := stats.New()
	c.IncTotal()
	c.IncTotal()
	if got := c.Summary().Total; got != 2 {
		t.Fatalf("expected Total=2, got %d", got)
	}
}

func TestCounter_IncMatched(t *testing.T) {
	c := stats.New()
	c.IncMatched()
	if got := c.Summary().Matched; got != 1 {
		t.Fatalf("expected Matched=1, got %d", got)
	}
}

func TestCounter_IncSkipped(t *testing.T) {
	c := stats.New()
	c.IncSkipped()
	c.IncSkipped()
	c.IncSkipped()
	if got := c.Summary().Skipped; got != 3 {
		t.Fatalf("expected Skipped=3, got %d", got)
	}
}

func TestCounter_IncErrored(t *testing.T) {
	c := stats.New()
	c.IncErrored()
	if got := c.Summary().Errored; got != 1 {
		t.Fatalf("expected Errored=1, got %d", got)
	}
}

func TestCounter_ElapsedIsPositive(t *testing.T) {
	c := stats.New()
	time.Sleep(2 * time.Millisecond)
	s := c.Summary()
	if s.Elapsed <= 0 {
		t.Fatalf("expected positive elapsed, got %v", s.Elapsed)
	}
}

func TestCounter_ConcurrentIncrements(t *testing.T) {
	c := stats.New()
	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func() {
			c.IncTotal()
			c.IncMatched()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 100; i++ {
		<-done
	}
	s := c.Summary()
	if s.Total != 100 {
		t.Fatalf("expected Total=100, got %d", s.Total)
	}
	if s.Matched != 100 {
		t.Fatalf("expected Matched=100, got %d", s.Matched)
	}
}
