package query

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
)

// CompiledQuery holds a parsed query with a compiled filter expression.
type CompiledQuery struct {
	Query
	Matcher *filter.Filter
}

// Compile takes a Query and compiles its filter expression into a reusable
// matcher. Returns an error if the filter expression is invalid.
func Compile(q Query) (*CompiledQuery, error) {
	if q.Filter == nil {
		return &CompiledQuery{Query: q, Matcher: nil}, nil
	}

	f, err := filter.Parse(*q.Filter)
	if err != nil {
		return nil, fmt.Errorf("compile filter: %w", err)
	}

	return &CompiledQuery{Query: q, Matcher: f}, nil
}

// Match reports whether the given log entry satisfies all constraints in the
// compiled query: time range and filter expression.
func (cq *CompiledQuery) Match(entry map[string]any) bool {
	if !matchTime(cq.Query, entry) {
		return false
	}
	if cq.Matcher != nil && !cq.Matcher.Match(entry) {
		return false
	}
	return true
}
