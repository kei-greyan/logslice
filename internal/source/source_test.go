package source_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"logslice/internal/source"
)

func TestFromArgs_EmptyUsesStdin(t *testing.T) {
	srcs, err := source.FromArgs([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(srcs) != 1 {
		t.Fatalf("expected 1 source, got %d", len(srcs))
	}
	if srcs[0].Type != source.TypeStdin {
		t.Errorf("expected TypeStdin, got %v", srcs[0].Type)
	}
	if srcs[0].Name != "stdin" {
		t.Errorf("expected name 'stdin', got %q", srcs[0].Name)
	}
}

func TestFromArgs_DashUsesStdin(t *testing.T) {
	srcs, err := source.FromArgs([]string{"-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srcs[0].Type != source.TypeStdin {
		t.Errorf("expected TypeStdin")
	}
}

func TestFromArgs_FileSource(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "app.log")
	if err := os.WriteFile(path, []byte(`{"msg":"hi"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	srcs, err := source.FromArgs([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(srcs) != 1 {
		t.Fatalf("expected 1 source, got %d", len(srcs))
	}
	if srcs[0].Type != source.TypeFile {
		t.Errorf("expected TypeFile")
	}
	if srcs[0].Name != "app.log" {
		t.Errorf("expected name 'app.log', got %q", srcs[0].Name)
	}
	data, _ := io.ReadAll(srcs[0].Reader)
	if string(data) != `{"msg":"hi"}` {
		t.Errorf("unexpected content: %s", data)
	}
	source.CloseAll(srcs)
}

func TestFromArgs_MissingFile(t *testing.T) {
	_, err := source.FromArgs([]string{"/nonexistent/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCloseAll_NoError(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "x.log")
	os.WriteFile(path, []byte("{}\n"), 0o644)

	srcs, _ := source.FromArgs([]string{path})
	if err := source.CloseAll(srcs); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}
