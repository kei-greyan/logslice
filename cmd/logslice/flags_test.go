package main

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.format != "pretty" {
		t.Errorf("expected default format 'pretty', got %q", cfg.format)
	}
	if cfg.filter != "" {
		t.Errorf("expected empty filter, got %q", cfg.filter)
	}
	if len(cfg.files) != 0 {
		t.Errorf("expected no files, got %v", cfg.files)
	}
}

func TestParseFlags_AllOptions(t *testing.T) {
	cfg, err := parseFlags([]string{
		"-filter", "level=error",
		"-since", "2024-01-01T00:00:00Z",
		"-until", "2024-12-31T23:59:59Z",
		"-format", "json",
		"app.log", "worker.log",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.filter != "level=error" {
		t.Errorf("expected filter 'level=error', got %q", cfg.filter)
	}
	if cfg.format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.format)
	}
	if len(cfg.files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.files))
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := parseFlags([]string{"-format", "xml"})
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestParseFlags_FormatCaseInsensitive(t *testing.T) {
	cfg, err := parseFlags([]string{"-format", "JSON"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.format != "json" {
		t.Errorf("expected normalized format 'json', got %q", cfg.format)
	}
}
