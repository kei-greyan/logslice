package ratelimit

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/reader"
)

func makeEntry(msg string) reader.Entry {
	return reader.Entry{
		Raw: []byte(`{"msg":"` + msg + `"}`),
		Fields: map[string]interface{}{"msg": msg},
	}
}

func feedEntries(entries []reader.Entry) <-chan reader.Entry {
	ch := make(chan reader.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func TestFilter_AllowsUnderBurst(t *testing.T) {
	entries := []reader.Entry{makeEntry("a"), makeEntry("b"), makeEntry("c")}
	in := feedEntries(entries)

	l := New(0, 5) // burst=5, no refill
	out := Filter(in, l, nil)

	var got []reader.Entry
	for e := range out {
		got = append(got, e)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestFilter_DropsWhenExhausted(t *testing.T) {
	entries := []reader.Entry{
		makeEntry("a"), makeEntry("b"), makeEntry("c"), makeEntry("d"),
	}
	in := feedEntries(entries)

	l := New(0, 2) // burst=2 only
	dropped := 0
	out := Filter(in, l, func() { dropped++ })

	var got []reader.Entry
	for e := range out {
		got = append(got, e)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if dropped != 2 {
		t.Fatalf("expected 2 dropped, got %d", dropped)
	}
}

func TestFilter_ClosesOutputWhenInputClosed(t *testing.T) {
	in := make(chan reader.Entry)
	close(in)
	l := New(10, 10)
	out := Filter(in, l, nil)

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}
