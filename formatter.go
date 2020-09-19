package logger

import "sort"

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
			s += "\""
			s += m[k]
			s += "\""

			if i+1 < len(keys) {
				s += ", "
			}
		}

		s += " }"
	}

	return s
}
