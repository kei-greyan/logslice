package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		input string
		field string
		op    string
		val   string
	}{
		{"level=error", "level", "=", "error"},
		{"level!=debug", "level", "!=", "debug"},
		{"msg~timeout", "msg", "~", "timeout"},
		{"status>=500", "status", ">=", "500"},
	}
	for _, tc := range cases {
		expr, err := filter.Parse(tc.input)
		if err != nil {
			t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
		}
		if expr.Field != tc.field || expr.Operator != tc.op || expr.Value != tc.val {
			t.Errorf("Parse(%q) = {%s %s %s}, want {%s %s %s}",
				tc.input, expr.Field, expr.Operator, expr.Value,
				tc.field, tc.op, tc.val)
		}
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := filter.Parse("noop")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestMatch(t *testing.T) {
	cases := []struct {
		expr  string
		line  string
		want  bool
	}{
		{"level=error", `{"level":"error","msg":"oops"}`, true},
		{"level=error", `{"level":"info","msg":"ok"}`, false},
		{"level!=debug", `{"level":"info"}`, true},
		{"msg~timeout", `{"msg":"connection timeout"}`, true},
		{"msg~timeout", `{"msg":"all good"}`, false},
		{"missing=x", `{"level":"info"}`, false},
	}
	for _, tc := range cases {
		expr, err := filter.Parse(tc.expr)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.expr, err)
		}
		got, err := expr.Match([]byte(tc.line))
		if err != nil {
			t.Fatalf("Match(%q): %v", tc.line, err)
		}
		if got != tc.want {
			t.Errorf("expr=%q line=%q: got %v, want %v", tc.expr, tc.line, got, tc.want)
		}
	}
}
