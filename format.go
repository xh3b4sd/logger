package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/xh3b4sd/logger/meta"
	"github.com/xh3b4sd/tracer"
)

var DefaultFormatter = JSONFormatter

func JSONIndenter(c context.Context, m map[string]string) string {
	str := JSONFormatter(c, m)

	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(str), "", "    ")
	if err != nil {
		tracer.Panic(err)
	}

	return buf.String()
}

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
	for k := range m {
		if k == KeyCal || k == KeyLev || k == KeyMes || k == KeySta || k == KeyTim {
			continue
		}

		if c != nil && meta.Has(c, k) {
			ctx = append(ctx, k)
		} else {
			oth = append(oth, k)
		}
	}

	if len(ctx) != 0 {
		slices.Sort(ctx)
	}

	{
		slices.Sort(oth)
	}

	var bui func(string) string
	{
		bui = func(k string) string {
			var s string
			{
				s += "\""
				s += k
				s += "\""
				s += ":"
			}

			if k != KeySta {
				s += "\""
			}

			{
				s += m[k]
			}

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
		s += "{ "
	}

	{
		s += bui(KeyTim)
		s += ", "
		s += bui(KeyLev)
		s += ", "
	}

	for _, k := range ctx {
		s += bui(k)
		s += ", "
	}

	{
		s += bui(KeyMes)
	}

	for i, k := range oth {
		if i == 0 {
			s += ", "
		}

		if m[k] == "" {
			continue
		}

		{
			s += bui(k)
		}

		if i+1 < len(oth) {
			s += ", "
		}
	}

	{
		_, e := m[KeySta]
		if e {
			if !strings.HasSuffix(s, ", ") {
				s += ", "
			}

			{
				s += bui(KeySta)
			}
		}
	}

	{
		_, e := m[KeyCal]
		if e {
			if !strings.HasSuffix(s, ", ") {
				s += ", "
			}

			{
				s += bui(KeyCal)
			}
		}
	}

	{
		s += " }"
	}

	return s
}
