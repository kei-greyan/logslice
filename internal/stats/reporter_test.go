package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"logslice/internal/stats"
)

func TestReporter_Print_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	r := stats.NewReporter(&buf)
	s := stats.Summary{
		Total:   50,
		Matched: 30,
		Skipped: 18,
		Errored: 2,
		Elapsed: 123 * time.Millisecond,
	}
	r.Print(s)

	out := buf.String()
	cases := []string{
		"logslice summary",
		"total   : 50",
		"matched : 30",
		"skipped : 18",
		"errored : 2",
		"elapsed",
	}
	for _, want := range cases {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestReporter_Print_ZeroSummary(t *testing.T) {
	var buf bytes.Buffer
	r := stats.NewReporter(&buf)
	r.Print(stats.Summary{})

	out := buf.String()
	if !strings.Contains(out, "total   : 0") {
		t.Errorf("expected zero total in output, got:\n%s", out)
	}
}

func TestReporter_Print_WritesToProvidedWriter(t *testing.T) {
	var buf bytes.Buffer
	r := stats.NewReporter(&buf)
	r.Print(stats.Summary{Total: 7})

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
