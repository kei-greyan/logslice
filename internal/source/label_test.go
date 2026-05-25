package source_test

import (
	"encoding/json"
	"testing"

	"logslice/internal/source"
)

func TestInjectSourceField_Adds(t *testing.T) {
	fields := map[string]any{"msg": "hello", "level": "info"}
	out := source.InjectSourceField(fields, "app.log", false)
	if out["_source"] != "app.log" {
		t.Errorf("expected _source=app.log, got %v", out["_source"])
	}
	if out["msg"] != "hello" {
		t.Errorf("original fields should be preserved")
	}
}

func TestInjectSourceField_NoOverwrite(t *testing.T) {
	fields := map[string]any{"_source": "original"}
	out := source.InjectSourceField(fields, "new", false)
	if out["_source"] != "original" {
		t.Errorf("expected _source to remain 'original', got %v", out["_source"])
	}
}

func TestInjectSourceField_Overwrite(t *testing.T) {
	fields := map[string]any{"_source": "original"}
	out := source.InjectSourceField(fields, "new", true)
	if out["_source"] != "new" {
		t.Errorf("expected _source='new', got %v", out["_source"])
	}
}

func TestInjectSourceField_DoesNotMutateInput(t *testing.T) {
	fields := map[string]any{"msg": "test"}
	source.InjectSourceField(fields, "x", false)
	if _, ok := fields["_source"]; ok {
		t.Error("original map should not be mutated")
	}
}

func TestMarshalLabeled_RoundTrip(t *testing.T) {
	e := source.LabeledEntry{
		Source: "svc.log",
		Fields: map[string]any{"level": "warn", "msg": "disk full"},
	}
	data, err := source.MarshalLabeled(e, false)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if out["_source"] != "svc.log" {
		t.Errorf("expected _source=svc.log, got %v", out["_source"])
	}
	if out["msg"] != "disk full" {
		t.Errorf("expected msg='disk full', got %v", out["msg"])
	}
}
