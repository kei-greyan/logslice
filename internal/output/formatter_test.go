package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func entry() map[string]any {
	return map[string]any{
		"time":  "2024-01-02T15:04:05Z",
		"level": "info",
		"msg":   "hello world",
		"user":  "alice",
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["msg"] != "hello world" {
		t.Errorf("expected msg=hello world, got %v", got["msg"])
	}
}

func TestFormatter_Pretty(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatPretty)

	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\n") {
		t.Error("expected indented (multi-line) output")
	}
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("pretty output is not valid JSON: %v", err)
	}
}

func TestFormatter_Text(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText)

	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "time=") {
		t.Errorf("expected text output to start with time=, got: %s", out)
	}
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected level=info in output: %s", out)
	}
	if !strings.Contains(out, "msg=hello world") {
		t.Errorf("expected msg=hello world in output: %s", out)
	}
	if !strings.Contains(out, "user=alice") {
		t.Errorf("expected user=alice in output: %s", out)
	}
}

func TestFormatter_TextPriorityOrder(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText)
	_ = f.Write(entry())

	out := strings.TrimSpace(buf.String())
	timeIdx := strings.Index(out, "time=")
	levelIdx := strings.Index(out, "level=")
	msgIdx := strings.Index(out, "msg=")

	if !(timeIdx < levelIdx && levelIdx < msgIdx) {
		t.Errorf("priority fields out of order: %s", out)
	}
}
