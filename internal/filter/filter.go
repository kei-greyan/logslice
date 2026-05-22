package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Expr represents a parsed filter expression.
type Expr struct {
	Field    string
	Operator string
	Value    string
}

// Parse parses a filter expression string of the form "field=value",
// "field!=value", "field~value" (contains), or "field>value" / "field<value".
func Parse(raw string) (*Expr, error) {
	ops := []string{"!=", ">=", "<=", "~", "=", ">", "<"}
	for _, op := range ops {
		if idx := strings.Index(raw, op); idx > 0 {
			return &Expr{
				Field:    strings.TrimSpace(raw[:idx]),
				Operator: op,
				Value:    strings.TrimSpace(raw[idx+len(op):]),
			}, nil
		}
	}
	return nil, fmt.Errorf("filter: invalid expression %q", raw)
}

// Match reports whether the JSON log line satisfies the expression.
func (e *Expr) Match(line []byte) (bool, error) {
	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return false, fmt.Errorf("filter: unmarshal: %w", err)
	}

	raw, ok := record[e.Field]
	if !ok {
		return false, nil
	}

	fieldVal := fmt.Sprintf("%v", raw)

	switch e.Operator {
	case "=":
		return fieldVal == e.Value, nil
	case "!=":
		return fieldVal != e.Value, nil
	case "~":
		return strings.Contains(fieldVal, e.Value), nil
	case ">":
		return fieldVal > e.Value, nil
	case "<":
		return fieldVal < e.Value, nil
	case ">=":
		return fieldVal >= e.Value, nil
	case "<=":
		return fieldVal <= e.Value, nil
	}
	return false, fmt.Errorf("filter: unknown operator %q", e.Operator)
}
