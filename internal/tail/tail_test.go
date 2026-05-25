package tail_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"logslice/internal/tail"
)

func writeLine(t *testing.T, f *os.File, entry map[string]any) {
	t.Helper()
	b, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	b = append(b, '\n')
	if _, err := f.Write(b); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tlr := tail.New(f.Name())
	if err := tlr.Start(); err != nil {
		t.Fatal(err)
	}
	defer tlr.Stop()

	// Allow the tailer goroutine to seek to EOF before we write.
	time.Sleep(50 * time.Millisecond)

	wantMsg := "hello from tail"
	writeLine(t, f, map[string]any{"msg": wantMsg, "level": "info"})

	select {
	case entry, ok := <-tlr.Out:
		if !ok {
			t.Fatal("channel closed unexpectedly")
		}
		if entry["msg"] != wantMsg {
			t.Errorf("got msg %q, want %q", entry["msg"], wantMsg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for tailed entry")
	}
}

func TestTailer_IgnoresPreExistingContent(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-pre-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Write a line BEFORE the tailer starts — it should be skipped.
	writeLine(t, f, map[string]any{"msg": "old line"})

	tlr := tail.New(f.Name())
	if err := tlr.Start(); err != nil {
		t.Fatal(err)
	}
	defer tlr.Stop()

	time.Sleep(50 * time.Millisecond)
	writeLine(t, f, map[string]any{"msg": "new line"})

	select {
	case entry := <-tlr.Out:
		if entry["msg"] != "new line" {
			t.Errorf("expected only new line, got %v", entry["msg"])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out")
	}
}

func TestTailer_StopClosesChannel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-stop-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	tlr := tail.New(f.Name())
	if err := tlr.Start(); err != nil {
		t.Fatal(err)
	}
	tlr.Stop()

	select {
	case _, ok := <-tlr.Out:
		if ok {
			t.Error("expected channel to be closed after Stop")
		}
	case <-time.After(time.Second):
		t.Fatal("channel was not closed after Stop")
	}
}
