package logger

import (
	"context"
	"sort"
	"strings"

	"github.com/xh3b4sd/logger/meta"
)

var DefaultFormatter = JSONFormatter

// JSONFormatter transforms the given map into the string representation of a
// simple JSON object. One specificity of JSONFormatter is that it treats the
// value of key "stack" differently. In case a key-value pair with the key
// "stack" is provided, the associated value is treated like a JSON object, and
// not like a JSON string. This difference in behaviour is unique compared to
// all other key-value pairs map may carry with it.
func JSONFormatter(c context.Context, m map[string]string) string {
	// Once we have the full map of key-value pairs we need the keys only in
	// order to sort them. We want the emitted log line to be alphabetically
	// ordered by keys.
	var ctx []string
	var oth []string
	{
		for k := range m {
			if k == KeyCal || k == KeyLev || k == KeyMes || k == KeySta || k == KeyTim {
				continue
			}

			if meta.Has(c, k) {
				ctx = append(ctx, k)
			} else {
				oth = append(oth, k)
			}
		}

		sort.Strings(ctx)
		sort.Strings(oth)
	}

	var key func(string) string
	{
		key = func(k string) string {
			if k == KeyCal || k == KeyCod || k == KeyDes || k == KeyDoc || k == KeyLev || k == KeyMes || k == KeySta || k == KeyTim {
				return k[:4]
			}

			return k
		}
	}

	var val func(string) string
	{
		val = func(k string) string {
			v, e := m[k]
			if !e {
				return ""
			}

			if k == KeyLev {
				return v[:4]
			}

			return v
		}
	}

	var bui func(string) string
	{
		bui = func(k string) string {
			var s string

			s += "\""
			s += key(k)
			s += "\""
			s += ":"

			if k != KeySta {
				s += "\""
			}

			s += val(k)

			if k != KeySta {
				s += "\""
			}

			return s
		}
	}

	// Below we compute a JSON string for the structured log line we want to
	// emit.
	var s string
	{
		{
			s += "{ "
		}

		{
			s += bui(KeyTim)
			s += ", "
			s += bui(KeyLev)
			s += ", "
		}

		{
			for _, k := range ctx {
				s += bui(k)
				s += ", "
			}
		}

		{
			s += bui(KeyMes)
		}

		{
			for i, k := range oth {
				if i == 0 {
					s += ", "
				}
				if m[k] == "" {
					continue
				}
				s += bui(k)
				if i+1 < len(oth) {
					s += ", "
				}
			}
		}

		{
			_, e := m[KeySta]
			if e {
				if !strings.HasSuffix(s, ", ") {
					s += ", "
				}
				s += bui(KeySta)
			}
		}

		{
			_, e := m[KeyCal]
			if e {
				if !strings.HasSuffix(s, ", ") {
					s += ", "
				}
				s += bui(KeyCal)
			}
		}

		{
			s += " }"
		}
	}

	return s
}
