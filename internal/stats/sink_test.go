package stats

import (
	"testing"
)

func TestSink_IncTotal(t *testing.T) {
	c := New()
	s := NewSink(c)
	s.IncTotal()
	s.IncTotal()
	sum := c.Summary()
	if sum.Total != 2 {
		t.Errorf("expected Total=2, got %d", sum.Total)
	}
}

func TestSink_IncMatched(t *testing.T) {
	c := New()
	s := NewSink(c)
	s.IncMatched()
	sum := c.Summary()
	if sum.Matched != 1 {
		t.Errorf("expected Matched=1, got %d", sum.Matched)
	}
}

func TestSink_IncSkipped(t *testing.T) {
	c := New()
	s := NewSink(c)
	s.IncSkipped()
	s.IncSkipped()
	s.IncSkipped()
	sum := c.Summary()
	if sum.Skipped != 3 {
		t.Errorf("expected Skipped=3, got %d", sum.Skipped)
	}
}

func TestSink_IncErrored(t *testing.T) {
	c := New()
	s := NewSink(c)
	s.IncErrored()
	sum := c.Summary()
	if sum.Errored != 1 {
		t.Errorf("expected Errored=1, got %d", sum.Errored)
	}
}

func TestSink_DoesNotShareStateWithOtherSinks(t *testing.T) {
	c1 := New()
	c2 := New()
	s1 := NewSink(c1)
	s2 := NewSink(c2)
	s1.IncTotal()
	s1.IncMatched()
	s2.IncTotal()
	s2.IncSkipped()
	sum1 := c1.Summary()
	sum2 := c2.Summary()
	if sum1.Matched != 1 || sum1.Skipped != 0 {
		t.Errorf("c1 unexpected: %+v", sum1)
	}
	if sum2.Matched != 0 || sum2.Skipped != 1 {
		t.Errorf("c2 unexpected: %+v", sum2)
	}
}

func TestSink_BackingCounterReflectsChanges(t *testing.T) {
	c := New()
	s := NewSink(c)
	for i := 0; i < 10; i++ {
		s.IncTotal()
	}
	if got := c.Summary().Total; got != 10 {
		t.Errorf("expected Total=10 via backing counter, got %d", got)
	}
}
