package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriter_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter("json", nil)
	if err != nil {
		t.Fatalf("NewFormatter: %v", err)
	}
	w := NewWriter(&buf, f)

	e := map[string]any{"level": "info", "msg": "hello"}
	if err := w.Write(e); err != nil {
		t.Fatalf("Write: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\"msg\":\"hello\"") && !strings.Contains(out, "\"msg\": \"hello\"") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("output should end with newline")
	}
}

func TestWriter_WriteAll(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter("json", nil)
	if err != nil {
		t.Fatalf("NewFormatter: %v", err)
	}
	w := NewWriter(&buf, f)

	entries := []map[string]any{
		{"msg": "first"},
		{"msg": "second"},
		{"msg": "third"},
	}
	if err := w.WriteAll(entries); err != nil {
		t.Fatalf("WriteAll: %v", err)
	}

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestWriter_WriteAll_Empty(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter("json", nil)
	if err != nil {
		t.Fatalf("NewFormatter: %v", err)
	}
	w := NewWriter(&buf, f)

	if err := w.WriteAll(nil); err != nil {
		t.Fatalf("WriteAll nil: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestWriter_ConcurrentWrite(t *testing.T) {
	var buf syncBuffer
	f, err := NewFormatter("json", nil)
	if err != nil {
		t.Fatalf("NewFormatter: %v", err)
	}
	w := NewWriter(&buf, f)

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(i int) {
			_ = w.Write(map[string]any{"i": i})
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

// syncBuffer is a thread-safe bytes.Buffer for testing.
type syncBuffer struct {
	mu  bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	return b.mu.Write(p)
}
