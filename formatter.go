package logger

import "sort"

// JSONFormatter transforms the given map into the string representation of a
// simple JSON object. One specificity of JSONFormatter is that it treats the
// value of key "stack" differently. In case a key-value pair with the key
// "stack" is provided, the associated value is treated like a JSON object, and
// not like a JSON string. This difference in behaviour is unique compared to
// all other key-value pairs map may carry with it.
func JSONFormatter(m map[string]string) string {
	// Once we have the full map of key-value pairs we need the keys only in
	// order to sort them. We want the emitted log line to be alphabetically
	// ordered by keys.
	var keys []string
	{
		for k := range m {
			keys = append(keys, k)
		}

		sort.Strings(keys)
	}

	// Below we compute a JSON string for the structured log line we want to
	// emit.
	var s string
	{
		s += "{ "

		for i, k := range keys {
			s += "\""
			s += k
			s += "\""
			s += ":"

			if k != KeyStack {
				s += "\""
			}

			s += m[k]

			if k != KeyStack {
				s += "\""
			}

			if i+1 < len(keys) {
				s += ", "
			}
		}

		s += " }"
	}

	return s
}
