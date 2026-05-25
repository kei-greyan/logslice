package source

import "encoding/json"

// LabeledEntry wraps a decoded JSON log entry with its origin source name.
type LabeledEntry struct {
	Source string
	Fields map[string]any
}

// InjectSourceField returns a new JSON object with the "_source" key set to
// the given label. If the original object already contains "_source" it is
// overwritten only when overwrite is true.
func InjectSourceField(raw map[string]any, label string, overwrite bool) map[string]any {
	if _, exists := raw["_source"]; exists && !overwrite {
		return raw
	}
	out := make(map[string]any, len(raw)+1)
	for k, v := range raw {
		out[k] = v
	}
	out["_source"] = label
	return out
}

// MarshalLabeled serialises a LabeledEntry back to compact JSON, injecting
// the source label into the output object.
func MarshalLabeled(e LabeledEntry, overwrite bool) ([]byte, error) {
	fields := InjectSourceField(e.Fields, e.Source, overwrite)
	return json.Marshal(fields)
}
