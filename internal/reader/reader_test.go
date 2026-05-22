package reader

import (
	"strings"
	"testing"
)

const sampleLogs = `{"level":"info","msg":"started"}
{"level":"error","msg":"failed","code":500}
not json at all
{"level":"debug","msg":"verbose"}
`

func TestReadReader_ParsesValidLines(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	records, errs := ReadReader("test", r)

	var got []Record
	for rec := range records {
		got = append(got, rec)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 records (non-JSON skipped), got %d", len(got))
	}
	for _, rec := range got {
		if rec.Source != "test" {
			t.Errorf("expected source %q, got %q", "test", rec.Source)
		}
		if rec.Data == nil {
			t.Error("expected non-nil Data map")
		}
		if rec.Raw == "" {
			t.Error("expected non-empty Raw field")
		}
	}
}

func TestReadReader_FieldValues(t *testing.T) {
	r := strings.NewReader(`{"level":"error","code":500}` + "\n")
	records, _ := ReadReader("src", r)

	rec := <-records
	if rec.Data["level"] != "error" {
		t.Errorf("expected level=error, got %v", rec.Data["level"])
	}
	if rec.Data["code"] != float64(500) {
		t.Errorf("expected code=500, got %v", rec.Data["code"])
	}
}

func TestReadReader_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	records, errs := ReadReader("empty", r)

	var count int
	for range records {
		count++
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 records, got %d", count)
	}
}

func TestReadReader_OnlyInvalidJSON(t *testing.T) {
	r := strings.NewReader("not json\nalso not json\n")
	records, errs := ReadReader("bad", r)

	var count int
	for range records {
		count++
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 records, got %d", count)
	}
}
