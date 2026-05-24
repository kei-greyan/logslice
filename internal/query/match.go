package query

import (
	"time"
)

// Entry is a parsed log entry with raw fields.
type Entry = map[string]interface{}

// Match reports whether a log entry satisfies the query constraints.
func (q *Query) Match(entry Entry) bool {
	if q.Filter != nil && !q.Filter.Match(entry) {
		return false
	}

	if q.Since != nil || q.Until != nil {
		t, ok := entryTime(entry, q.TimeField)
		if !ok {
			// If we can't parse the time, exclude the entry when time bounds are set.
			return false
		}
		if q.Since != nil && t.Before(*q.Since) {
			return false
		}
		if q.Until != nil && t.After(*q.Until) {
			return false
		}
	}

	return true
}

// entryTime extracts and parses a time.Time from the named field of an entry.
func entryTime(entry Entry, field string) (time.Time, bool) {
	v, ok := entry[field]
	if !ok {
		return time.Time{}, false
	}

	switch val := v.(type) {
	case string:
		for _, layout := range timeLayouts {
			if t, err := time.Parse(layout, val); err == nil {
				return t, true
			}
		}
	case float64:
		// Unix timestamp in seconds (JSON numbers).
		return time.Unix(int64(val), 0).UTC(), true
	}

	return time.Time{}, false
}
